package authorizationsync

import (
	"fmt"

	"github.com/golang/glog"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/rbac"
	rbacclient "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/typed/rbac/internalversion"
	rbacinformers "k8s.io/kubernetes/pkg/client/informers/informers_generated/internalversion/rbac/internalversion"
	rbaclister "k8s.io/kubernetes/pkg/client/listers/rbac/internalversion"
	"k8s.io/kubernetes/pkg/controller"

	authorizationapi "github.com/openshift/origin/pkg/authorization/api"
	origininformers "github.com/openshift/origin/pkg/authorization/generated/informers/internalversion/authorization/internalversion"
	originlister "github.com/openshift/origin/pkg/authorization/generated/listers/authorization/internalversion"
)

type OriginClusterRoleBindingToRBACClusterRoleBindingController struct {
	rbacClient rbacclient.ClusterRoleBindingsGetter

	rbacLister    rbaclister.ClusterRoleBindingLister
	originIndexer cache.Indexer
	originLister  originlister.ClusterRoleBindingLister

	genericController
}

func NewOriginToRBACClusterRoleBindingController(rbacClusterRoleBindingInformer rbacinformers.ClusterRoleBindingInformer, originClusterPolicyBindingInformer origininformers.ClusterPolicyBindingInformer, rbacClient rbacclient.ClusterRoleBindingsGetter) *OriginClusterRoleBindingToRBACClusterRoleBindingController {
	originIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
	c := &OriginClusterRoleBindingToRBACClusterRoleBindingController{
		rbacClient:    rbacClient,
		rbacLister:    rbacClusterRoleBindingInformer.Lister(),
		originIndexer: originIndexer,
		originLister:  originlister.NewClusterRoleBindingLister(originIndexer),

		genericController: genericController{
			name: "OriginClusterRoleBindingToRBACClusterRoleBindingController",
			cachesSynced: func() bool {
				return rbacClusterRoleBindingInformer.Informer().HasSynced() && originClusterPolicyBindingInformer.Informer().HasSynced()
			},
			queue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "origin-to-rbac-rolebinding"),
		},
	}
	c.genericController.syncFunc = c.syncClusterRoleBinding

	rbacClusterRoleBindingInformer.Informer().AddEventHandler(naiveEventHandler(c.queue))
	originClusterPolicyBindingInformer.Informer().AddEventHandler(c.clusterPolicyBindingEventHandler())

	return c
}

func (c *OriginClusterRoleBindingToRBACClusterRoleBindingController) syncClusterRoleBinding(name string) error {
	rbacClusterRoleBinding, rbacErr := c.rbacLister.Get(name)
	if !apierrors.IsNotFound(rbacErr) && rbacErr != nil {
		return rbacErr
	}
	originClusterRoleBinding, originErr := c.originLister.Get(name)
	if !apierrors.IsNotFound(originErr) && originErr != nil {
		return originErr
	}

	// if neither roleBinding exists, return
	if apierrors.IsNotFound(rbacErr) && apierrors.IsNotFound(originErr) {
		return nil
	}
	// if the origin roleBinding doesn't exist, just delete the rbac roleBinding
	if apierrors.IsNotFound(originErr) {
		// orphan on delete to minimize fanout.  We ought to clean the rest via controller too.
		return c.rbacClient.ClusterRoleBindings().Delete(name, nil)
	}

	// convert the origin roleBinding to an rbac roleBinding and compare the results
	convertedClusterRoleBinding := &rbac.ClusterRoleBinding{}
	if err := authorizationapi.Convert_api_ClusterRoleBinding_To_rbac_ClusterRoleBinding(originClusterRoleBinding, convertedClusterRoleBinding, nil); err != nil {
		return err
	}
	// do a deep copy here since conversion does not guarantee a new object.
	equivalentClusterRoleBinding := &rbac.ClusterRoleBinding{}
	if err := rbac.DeepCopy_rbac_ClusterRoleBinding(convertedClusterRoleBinding, equivalentClusterRoleBinding, cloner); err != nil {
		return err
	}

	// if we're missing the rbacClusterRoleBinding, create it
	if apierrors.IsNotFound(rbacErr) {
		equivalentClusterRoleBinding.ResourceVersion = ""
		_, err := c.rbacClient.ClusterRoleBindings().Create(equivalentClusterRoleBinding)
		return err
	}

	// if we might need to update, we need to stomp fields that are never going to match like uid and creation time
	equivalentClusterRoleBinding.SelfLink = rbacClusterRoleBinding.SelfLink
	equivalentClusterRoleBinding.UID = rbacClusterRoleBinding.UID
	equivalentClusterRoleBinding.ResourceVersion = rbacClusterRoleBinding.ResourceVersion
	equivalentClusterRoleBinding.CreationTimestamp = rbacClusterRoleBinding.CreationTimestamp

	// if they're equal, we have no work to do
	if kapi.Semantic.DeepEqual(equivalentClusterRoleBinding, rbacClusterRoleBinding) {
		return nil
	}

	glog.V(1).Infof("writing RBAC clusterrolebinding %v", name)
	_, err := c.rbacClient.ClusterRoleBindings().Update(equivalentClusterRoleBinding)
	// if the update was invalid, we're probably changing an immutable field or something like that
	// either way, the existing object is wrong.  Delete it and try again.
	if apierrors.IsInvalid(err) {
		c.rbacClient.ClusterRoleBindings().Delete(name, nil)
	}
	return err
}

func (c *OriginClusterRoleBindingToRBACClusterRoleBindingController) clusterPolicyBindingEventHandler() cache.ResourceEventHandler {
	return cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			originContainerObj := obj.(*authorizationapi.ClusterPolicyBinding)
			for _, originObj := range originContainerObj.RoleBindings {
				c.originIndexer.Add(originObj)
				key, err := controller.KeyFunc(originObj)
				if err != nil {
					utilruntime.HandleError(err)
					continue
				}
				c.queue.Add(key)
			}
		},
		UpdateFunc: func(old, cur interface{}) {
			originContainerObj := cur.(*authorizationapi.ClusterPolicyBinding)
			for _, originObj := range originContainerObj.RoleBindings {
				c.originIndexer.Add(originObj)
				key, err := controller.KeyFunc(originObj)
				if err != nil {
					utilruntime.HandleError(err)
					continue
				}
				c.queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			originContainerObj, ok := obj.(*authorizationapi.ClusterPolicyBinding)
			if !ok {
				tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
				if !ok {
					utilruntime.HandleError(fmt.Errorf("Couldn't get object from tombstone %#v", obj))
				}
				originContainerObj, ok = tombstone.Obj.(*authorizationapi.ClusterPolicyBinding)
				if !ok {
					utilruntime.HandleError(fmt.Errorf("Tombstone contained object that is not a runtime.Object %#v", obj))
				}
			}

			for _, originObj := range originContainerObj.RoleBindings {
				c.originIndexer.Add(originObj)
				key, err := controller.KeyFunc(originObj)
				if err != nil {
					utilruntime.HandleError(err)
					continue
				}
				c.queue.Add(key)
			}
		},
	}
}
