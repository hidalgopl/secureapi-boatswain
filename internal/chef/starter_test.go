package chef
//
//import (
//	"context"
//	"testing"
//
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//	"go.sauce.io/ondemand/boatswain/internal/testtools/goldenfile"
//	"go.sauce.io/ondemand/boatswain/internal/testtools/testjson"
//
//	"k8s.io/api/core/v1"
//)
//
//func TestStartChef(t *testing.T) {
//	fakeCreator := newFakeCreator()
//	s := defaultStarter{
//		podCreator:  fakeCreator,
//		uuidGenFunc: fakeUUIDGenFunc,
//		options: Options{
//			ChefImage:        "chef-image",
//			ChefTag:          "chef-tag",
//			BrowserSyncImage: "browser-sync-image",
//		},
//	}
//
//	req := startChefReq{
//		Browser: browser{
//			Name:    "chrome",
//			Version: "72.0.3626.121",
//		},
//	}
//
//	rsp, err := s.start(context.TODO(), req)
//	assert.NoError(t, err)
//
//	// run:
//	// go test go.sauce.io/ondemand/boatswain/internal/chef -test.update -run TestStartChef
//	// to update chef_req.golden
//
//	actualReq := fakeCreator.reqPodSpecAsJSON(t)
//	expectedReq := goldenfile.Get(t, actualReq, "chef_req.golden")
//	assert.JSONEq(t, expectedReq, actualReq)
//
//	assert.Equal(t, 0, rsp.Status)
//	assert.Equal(t, "chef-fake-uuid", rsp.PodStatus.Name)
//	assert.Equal(t, "Running", rsp.PodStatus.Phase)
//	assert.Equal(t, "fake-ip", rsp.PodStatus.IP)
//}
//
//func TestStartChef_Select_Browser_Dir(t *testing.T) {
//	tt := []struct {
//		name        string
//		browser     string
//		version     string
//		expectedDir string
//		skip        bool
//	}{
//		{
//			name:        "chrome full version",
//			browser:     "chrome",
//			version:     "72.0.3626.121",
//			expectedDir: "/browsers/chrome-72",
//		},
//		{
//			name:        "chrome major only",
//			browser:     "chrome",
//			version:     "72",
//			expectedDir: "/browsers/chrome-72",
//		},
//		{
//			name:        "firefox full version",
//			browser:     "firefox",
//			version:     "65.0.2",
//			expectedDir: "/browsers/firefox-65.0",
//		},
//		//TODO(core) gather info abut supported browsers & versions then add missing tests & refactor code
//		{
//			name:        "firefox full version",
//			browser:     "firefox",
//			version:     "65",
//			expectedDir: "/browsers/firefox-65",
//			skip:        true,
//		},
//	}
//
//	for _, tc := range tt {
//		t.Run(tc.name, func(t *testing.T) {
//			if tc.skip {
//				t.SkipNow()
//			}
//
//			fakeCreator := newFakeCreator()
//			service := defaultStarter{
//				podCreator:  fakeCreator,
//				uuidGenFunc: fakeUUIDGenFunc,
//			}
//
//			req := startChefReq{
//				Browser: browser{
//					Name:    tc.browser,
//					Version: tc.version,
//				},
//			}
//			_, err := service.start(context.TODO(), req)
//			require.NoError(t, err)
//
//			actualCmd := fakeCreator.reqPodSpec.Spec.Containers[0].ReadinessProbe.Exec.Command
//			expectedCmd := []string{"test", "-f", "/tmp/ready", "-a", "-d", tc.expectedDir}
//			assert.Equal(t, expectedCmd, actualCmd)
//		})
//	}
//}
//
//func TestStartChef_NodeSelector_In_Local_Mode(t *testing.T) {
//	fakeCreator := newFakeCreator()
//	s := defaultStarter{
//		podCreator:  fakeCreator,
//		uuidGenFunc: fakeUUIDGenFunc,
//		options: Options{
//			LocalMode: false,
//		},
//	}
//
//	req := startChefReq{
//		Browser: browser{
//			Name:    "chrome",
//			Version: "72",
//		},
//	}
//
//	_, err := s.start(context.TODO(), req)
//	assert.NoError(t, err)
//
//	expectedNodeSelector := map[string]string{
//		"chefonly": "true",
//	}
//
//	assert.Equal(t, expectedNodeSelector, fakeCreator.reqPodSpec.Spec.NodeSelector)
//}
//
//func TestStartHeadlessChef(t *testing.T) {
//	fakeCreator := newFakeCreator()
//	s := defaultStarter{
//		podCreator:  fakeCreator,
//		uuidGenFunc: fakeUUIDGenFunc,
//		options: Options{
//			ChefDryads:        "dryad:latest",
//			HeadlessChefImage: "chefImage",
//			HeadlessChefTag:   "chefTag",
//		},
//	}
//
//	req := startChefReq{
//		Browser: browser{
//			Name:    "chrome",
//			Version: "72.0.3626.121",
//		},
//	}
//
//	rsp, err := s.startHeadless(context.TODO(), req)
//	assert.NoError(t, err)
//
//	// run:
//	// go test go.sauce.io/ondemand/boatswain/internal/chef -test.update -run TestStartHeadlessChef
//	// to update headless_chef_req.golden
//
//	actualPodReq := fakeCreator.reqPodSpecAsJSON(t)
//	expectedPodReq := goldenfile.Get(t, actualPodReq, "headless_chef_req.golden")
//	assert.JSONEq(t, expectedPodReq, actualPodReq)
//
//	assert.Equal(t, 0, rsp.Status)
//	assert.Equal(t, "chef-fake-uuid", rsp.PodStatus.Name)
//	assert.Equal(t, "Running", rsp.PodStatus.Phase)
//	assert.Equal(t, "fake-ip", rsp.PodStatus.IP)
//}
//
//func TestStartHeadlessChef_NodeSelector_In_Local_Mode(t *testing.T) {
//	fakeCreator := newFakeCreator()
//	s := defaultStarter{
//		podCreator:  fakeCreator,
//		uuidGenFunc: fakeUUIDGenFunc,
//		options: Options{
//			LocalMode: false,
//		},
//	}
//
//	req := startChefReq{
//		Browser: browser{
//			Name:    "chrome",
//			Version: "72",
//		},
//	}
//
//	_, err := s.startHeadless(context.TODO(), req)
//	assert.NoError(t, err)
//
//	expectedNodeSelector := map[string]string{
//		"chefonly": "true",
//	}
//
//	assert.Equal(t, expectedNodeSelector, fakeCreator.reqPodSpec.Spec.NodeSelector)
//}
//
//type fakeCreator struct {
//	reqPodSpec *v1.Pod
//}
//
//func newFakeCreator() *fakeCreator {
//	return &fakeCreator{}
//}
//
//func (fc *fakeCreator) Create(ctx context.Context, podSpec *v1.Pod) (*v1.Pod, error) {
//	fc.reqPodSpec = podSpec
//	runningPod := podSpec.DeepCopy()
//	runningPod.Status = v1.PodStatus{
//		Phase: v1.PodRunning,
//		PodIP: "fake-ip",
//		Conditions: []v1.PodCondition{
//			{Type: v1.PodReady, Status: v1.ConditionTrue},
//		},
//	}
//	return runningPod, nil
//}
//
//func (fc *fakeCreator) reqPodSpecAsJSON(t *testing.T) string {
//	return testjson.ToJSON(t, fc.reqPodSpec)
//}
//
//func fakeUUIDGenFunc() string {
//	return "fake-uuid"
//}
