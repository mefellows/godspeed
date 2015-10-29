package strategy

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/mefellows/godspeed/log"
	"github.com/mefellows/plugo/plugo"
	"math/rand"
	"time"
)

type ECSDeploymentStrategy struct {
	ClusterName    string `mapstructure:"cluster_name"`
	Application    string
	Region         string `default:"ap-southeast-2"`
	TaskDefinition string `mapstructure:"task_definition"`
	ecs            *ecs.ECS
	ELB            string
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	plugo.PluginFactories.Register(func() (interface{}, error) {
		return &ECSDeploymentStrategy{}, nil
	}, "ecs")
}

func (s *ECSDeploymentStrategy) Setup() {
	log.Debug("Setting up ECS ")

	s.ecs = ecs.New(nil)
}

func (s *ECSDeploymentStrategy) createCommand(command string) {
}

func (s *ECSDeploymentStrategy) checkCluster(cluster string) {

	params := &ecs.DescribeClustersInput{
		Clusters: []*string{
			aws.String(cluster),
		},
	}
	resp, err := s.ecs.DescribeClusters(params)

	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	if len(resp.Clusters) == 0 {
		log.Fatalf("Cluster '%s' does not exist", cluster)
	}

	if *resp.Clusters[0].Status != "ACTIVE" {
		log.Info("Cluster not yet active, waiting...")
	}
}

func (s *ECSDeploymentStrategy) Deploy() error {
	log.Info("Deploying to ECS")

	s.checkCluster(s.ClusterName)
	//s.registerTask()

	return nil
}

func (s *ECSDeploymentStrategy) Rollback() error {
	log.Info("ECS Rolling back")
	return nil
}

func (s *ECSDeploymentStrategy) Teardown() {
	log.Debug("ECS Teardown")
}
