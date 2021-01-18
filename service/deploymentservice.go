/**
 * @Author: lzw5399
 * @Date: 2021/1/17 16:39
 * @Desc:
 */
package service

type DeploymentService interface {
	CreateDeployment()
}

type deploymentService struct {
}

func NewDeploymentService() DeploymentService {
	return deploymentService{}
}

func (d deploymentService) CreateDeployment() {

}
