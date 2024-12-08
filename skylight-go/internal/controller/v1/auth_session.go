package v1

import (
	"skylight/internal/model/entity"
	"skylight/internal/service"
	os "skylight/internal/service/openstack"
	"strings"

	"github.com/BytemanD/skyman/openstack"
	"github.com/BytemanD/skyman/openstack/auth"
	"github.com/gogf/gf/v2/net/ghttp"
)

type AuthSession struct {
	AuthUrl     string
	UserName    string
	ProjectName string
	Password    string

	webSocket       *ghttp.WebSocket
	OpenstackClient *openstack.Openstack
}

func (s *AuthSession) SetRegion(region string) {
	s.OpenstackClient.AuthPlugin = auth.NewPasswordAuth(
		s.AuthUrl, auth.User{Name: s.UserName},
		auth.Project{Name: s.ProjectName},
		region,
	)
}

func (s *AuthSession) UpdateSocket(req *ghttp.Request) error {
	ws, err := req.WebSocket()
	if err != nil {
		return err
	}
	s.webSocket = ws
	return nil
}

func (s *AuthSession) Publish(message entity.Message) error {
	return s.webSocket.WriteJSON(message)
}

var Session map[string]*AuthSession

func GetAuthSession(req *ghttp.Request) *AuthSession {
	session := Session[req.GetSessionId()]
	if session == nil {
		sessionLoginInfo, _ := req.Session.Get("loginInfo", nil)
		loginInfo := os.LoginInfo{}
		if err := sessionLoginInfo.Struct(&loginInfo); err != nil {
			return nil
		}
		cluster, _ := service.ClusterService.GetClusterByName(loginInfo.Cluster)
		if cluster == nil {
			return nil
		}
		AddAuthSession(
			req.GetSessionId(),
			cluster.AuthUrl, loginInfo.User.Name, loginInfo.Project.Name,
			loginInfo.Password, loginInfo.Region,
		)
	}
	return Session[req.GetSessionId()]
}

func AddAuthSession(sessionId string, authUrl, username, projectName, password, region string) {
	if !strings.HasSuffix(authUrl, "/v3") {
		authUrl += "/v3"
	}
	Session[sessionId] = &AuthSession{
		AuthUrl:     authUrl,
		UserName:    username,
		ProjectName: projectName,
		OpenstackClient: openstack.NewClient(
			authUrl, auth.User{Name: username, Password: password, Domain: auth.Domain{Name: "Default"}},
			auth.Project{Name: projectName, Domain: auth.Domain{Name: "Default"}},
			region,
		),
	}
}

func init() {
	Session = map[string]*AuthSession{}
}
