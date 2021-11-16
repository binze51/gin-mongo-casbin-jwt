package manager

import (
	"context"

	"gitee.com/binze/binzekeji/model/manager"
	svc "gitee.com/binze/binzekeji/service/manager"
	"go.mongodb.org/mongo-driver/bson"
)

//根据根角色找其所有子角色，递归成多节点树
func findChildrenRole(parent *manager.SysRole) (err error) {
	ctx := context.TODO()
	filter := bson.M{"parentId": parent.ID}
	curr := svc.ManagerSvc.Collection(manager.TSys_roles).Find(ctx, filter)
	curr.All(&parent.Children)
	if len(parent.Children) > 0 {
		for k := range parent.Children {
			err = findChildrenRole(&parent.Children[k])
		}
	}
	return err
}

//根据根菜单找其所有子菜单，递归成多节点树
func findChildrenMenu(parent *manager.SysMenu) (err error) {
	ctx := context.TODO()
	filter := bson.M{"parentId": parent.ID}
	curr := svc.ManagerSvc.Collection(manager.TSys_roles).Find(ctx, filter)
	curr.All(&parent.Children)
	if len(parent.Children) > 0 {
		for k := range parent.Children {
			err = findChildrenMenu(&parent.Children[k])
		}
	}
	return err
}
