package v1

import (
	"context"
	"sync"

	"github.com/neee333ko/IAM/internal/apiserver/store"
	v1 "github.com/neee333ko/api/apiserver/v1"
	metav1 "github.com/neee333ko/component-base/pkg/meta/v1"
)

type UserServ interface {
	Create(c context.Context, user *v1.User, options *metav1.CreateOptions) error
	Update(c context.Context, user *v1.User, options *metav1.UpdateOptions) error
	Get(c context.Context, name string, options *metav1.GetOptions) (*v1.User, error)
	List(c context.Context, options *metav1.ListOptions) (*v1.UserList, error)
	ListWithBadPerformance(c context.Context, options *metav1.ListOptions) (*v1.UserList, error)
	Delete(c context.Context, name string, options *metav1.DeleteOptions) error
	DeleteCollection(c context.Context, names []string, options *metav1.DeleteOptions) error
}

type UserService struct {
	store store.Factory
}

func (us *UserService) Create(c context.Context, user *v1.User, options *metav1.CreateOptions) error {
	if err := us.store.NewUserStore().Create(c, user, options); err != nil {
		return err
	}

	return nil
}

func (us *UserService) Update(c context.Context, user *v1.User, options *metav1.UpdateOptions) error {
	if err := us.store.NewUserStore().Update(c, user, options); err != nil {
		return err
	}

	return nil
}

func (us *UserService) Get(c context.Context, name string, options *metav1.GetOptions) (*v1.User, error) {
	user, err := us.store.NewUserStore().Get(c, name, options)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) List(c context.Context, options *metav1.ListOptions) (*v1.UserList, error) {
	userlist, err := us.store.NewUserStore().List(c, options)
	if err != nil {
		return nil, err
	}

	errchan := make(chan error)
	finished := make(chan int)

	var wg sync.WaitGroup
	var newitems sync.Map

	for _, user := range userlist.Items {
		wg.Add(1)

		go func(user *v1.User) {
			defer wg.Done()

			pl, err := us.store.NewPolicyStore().Get(c, user.Name, &metav1.GetOptions{})
			if err != nil {
				errchan <- err
				return
			}

			newitems.Store(user.ID, &v1.User{
				ObjectMeta: metav1.ObjectMeta{
					ID:         user.ID,
					InstanceID: user.InstanceID,
					Name:       user.Name,
					Ext:        user.Ext,
					CreatedAt:  user.CreatedAt,
					UpdatedAt:  user.UpdatedAt,
				},
				Status:      user.Status,
				Nickname:    user.Nickname,
				Password:    user.Password,
				Email:       user.Email,
				Phone:       user.Phone,
				IsAdmin:     user.IsAdmin,
				TotalPolicy: pl.GetTotalCount(),
				LoginedAt:   user.LoginedAt,
			})
		}(user)
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case err = <-errchan:
		return nil, err
	case <-finished:
		for i, user := range userlist.Items {
			u, _ := newitems.Load(user.ID)
			userlist.Items[i] = u.(*v1.User)
		}

		return userlist, nil
	}
}

func (us *UserService) ListWithBadPerformance(c context.Context, options *metav1.ListOptions) (*v1.UserList, error) {
	userlist, err := us.store.NewUserStore().List(c, options)
	if err != nil {
		return nil, err
	}

	for i, user := range userlist.Items {
		pl, err := us.store.NewPolicyStore().Get(c, user.Name, &metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		userlist.Items[i].TotalPolicy = pl.GetTotalCount()
	}

	return userlist, nil
}

func (us *UserService) Delete(c context.Context, name string, options *metav1.DeleteOptions) error {
	if err := us.store.NewUserStore().Delete(c, name, options); err != nil {
		return err
	}

	return nil
}

func (us *UserService) DeleteCollection(c context.Context, names []string, options *metav1.DeleteOptions) error {
	if err := us.store.NewUserStore().DeleteCollection(c, names, options); err != nil {
		return err
	}

	return nil
}
