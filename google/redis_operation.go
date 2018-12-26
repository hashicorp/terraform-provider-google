package google

import (
	"google.golang.org/api/redis/v1beta1"
)

type RedisOperationWaiter struct {
	Service *redis.ProjectsLocationsService
	CommonOperationWaiter
}

func (w *RedisOperationWaiter) QueryOp() (interface{}, error) {
	return w.Service.Operations.Get(w.Op.Name).Do()
}

func redisOperationWaitTime(service *redis.Service, op *redis.Operation, project, activity string, timeoutMinutes int) error {
	w := &RedisOperationWaiter{
		Service: service.Projects.Locations,
	}
	if err := w.SetOp(op); err != nil {
		return err
	}
	return OperationWait(w, activity, timeoutMinutes)
}
