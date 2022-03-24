package manager

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"

	"github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

// UpdateStackOptions contains options for updating a stack.
type UpdateStackOptions struct {
	Stack         *Stack
	StackName     string
	ChangeSetName string
	Description   string
	TemplateData  TemplateData
	Parameters    map[string]string
	Wait          bool
}

// GetNodegroupOption nodegroup options.
type GetNodegroupOption struct {
	Stack         *NodeGroupStack
	NodeGroupName string
}

var _ StackManager = &StackCollection{}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_stack_manager.go . StackManager
type StackManager interface {
	AppendNewClusterStackResource(ctx context.Context, plan, supportsManagedNodes bool) (bool, error)
	CreateStack(ctx context.Context, name string, stack builder.ResourceSet, tags, parameters map[string]string, errs chan error) error
	DeleteStackBySpec(ctx context.Context, s *Stack) (*Stack, error)
	DeleteStackBySpecSync(ctx context.Context, s *Stack, errs chan error) error
	DeleteStackSync(ctx context.Context, s *Stack) error
	DeleteTasksForDeprecatedStacks(ctx context.Context) (*tasks.TaskTree, error)
	DescribeClusterStack(ctx context.Context) (*Stack, error)
	DescribeIAMServiceAccountStacks(ctx context.Context) ([]*Stack, error)
	DescribeNodeGroupStack(ctx context.Context, nodeGroupName string) (*Stack, error)
	DescribeNodeGroupStacks(ctx context.Context) ([]*Stack, error)
	DescribeNodeGroupStacksAndResources(ctx context.Context) (map[string]StackInfo, error)
	DescribeStack(ctx context.Context, i *Stack) (*Stack, error)
	DescribeStackChangeSet(ctx context.Context, i *Stack, changeSetName string) (*ChangeSet, error)
	DescribeStackEvents(ctx context.Context, i *Stack) ([]types.StackEvent, error)
	DescribeStacks(ctx context.Context) ([]*Stack, error)
	DoCreateStackRequest(ctx context.Context, i *Stack, templateData TemplateData, tags, parameters map[string]string, withIAM bool, withNamedIAM bool) error
	DoWaitUntilStackIsCreated(i *Stack) error
	EnsureMapPublicIPOnLaunchEnabled(ctx context.Context) error
	FixClusterCompatibility(ctx context.Context) error
	GetAutoScalingGroupDesiredCapacity(ctx context.Context, name string) (autoscaling.Group, error)
	GetAutoScalingGroupName(ctx context.Context, s *Stack) (string, error)
	GetClusterStackIfExists(ctx context.Context) (*Stack, error)
	GetFargateStack(ctx context.Context) (*Stack, error)
	GetIAMAddonName(s *Stack) string
	GetIAMAddonsStacks(ctx context.Context) ([]*Stack, error)
	GetIAMServiceAccounts(ctx context.Context) ([]*v1alpha5.ClusterIAMServiceAccount, error)
	GetKarpenterStack(ctx context.Context) (*Stack, error)
	GetManagedNodeGroupTemplate(ctx context.Context, options GetNodegroupOption) (string, error)
	GetNodeGroupName(s *Stack) string
	GetNodeGroupStackType(ctx context.Context, options GetNodegroupOption) (v1alpha5.NodeGroupType, error)
	GetStackTemplate(ctx context.Context, stackName string) (string, error)
	GetUnmanagedNodeGroupAutoScalingGroupName(ctx context.Context, s *Stack) (string, error)
	HasClusterStackFromList(ctx context.Context, clusterStackNames []string, clusterName string) (bool, error)
	ListClusterStackNames(ctx context.Context) ([]string, error)
	ListIAMServiceAccountStacks(ctx context.Context) ([]string, error)
	ListNodeGroupStacks(ctx context.Context) ([]NodeGroupStack, error)
	ListStacks(ctx context.Context, statusFilters ...types.StackStatus) ([]*Stack, error)
	ListStacksMatching(ctx context.Context, nameRegex string, statusFilters ...types.StackStatus) ([]*Stack, error)
	LookupCloudTrailEvents(ctx context.Context, i *Stack) ([]*cloudtrail.Event, error)
	MakeChangeSetName(action string) string
	MakeClusterStackName() string
	NewClusterCompatTask() tasks.Task
	NewManagedNodeGroupTask(nodeGroups []*v1alpha5.ManagedNodeGroup, forceAddCNIPolicy bool, importer vpc.Importer) *tasks.TaskTree
	NewTaskToDeleteAddonIAM(ctx context.Context, wait bool) (*tasks.TaskTree, error)
	NewTaskToDeleteUnownedNodeGroup(clusterName, nodegroup string, eksAPI eksiface.EKSAPI, waitCondition *DeleteWaitCondition) tasks.Task
	NewTasksToCreateClusterWithNodeGroups(nodeGroups []*v1alpha5.NodeGroup, managedNodeGroups []*v1alpha5.ManagedNodeGroup, postClusterCreationTasks ...tasks.Task) *tasks.TaskTree
	NewTasksToCreateIAMServiceAccounts(serviceAccounts []*v1alpha5.ClusterIAMServiceAccount, oidc *iamoidc.OpenIDConnectManager, clientSetGetter kubernetes.ClientSetGetter) *tasks.TaskTree
	NewTasksToDeleteClusterWithNodeGroups(ctx context.Context, stack *Stack, stacks []NodeGroupStack, deleteOIDCProvider bool, oidc *iamoidc.OpenIDConnectManager, clientSetGetter kubernetes.ClientSetGetter, wait bool, cleanup func(chan error, string) error) (*tasks.TaskTree, error)
	NewTasksToDeleteIAMServiceAccounts(ctx context.Context, serviceAccounts []string, clientSetGetter kubernetes.ClientSetGetter, wait bool) (*tasks.TaskTree, error)
	NewTasksToDeleteNodeGroups(stacks []NodeGroupStack, shouldDelete func(_ string) bool, wait bool, cleanup func(chan error, string) error) (*tasks.TaskTree, error)
	NewTasksToDeleteOIDCProviderWithIAMServiceAccounts(ctx context.Context, oidc *iamoidc.OpenIDConnectManager, clientSetGetter kubernetes.ClientSetGetter) (*tasks.TaskTree, error)
	NewUnmanagedNodeGroupTask(nodeGroups []*v1alpha5.NodeGroup, forceAddCNIPolicy bool, importer vpc.Importer) *tasks.TaskTree
	RefreshFargatePodExecutionRoleARN(ctx context.Context) error
	StackStatusIsNotReady(s *Stack) bool
	StackStatusIsNotTransitional(s *Stack) bool
	UpdateNodeGroupStack(ctx context.Context, nodeGroupName, template string, wait bool) error
	UpdateStack(ctx context.Context, options UpdateStackOptions) error
}
