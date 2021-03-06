package deploy

import (
	"errors"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/context"
)

type deployS struct {
	Project  string  `json:"project,omitempty"`
	Name     string  `json:"name,omitempty"`
	Template string  `json:"template,omitempty"`
	Image    string  `json:"image,omitempty"`
	Url      string  `json:"url,omitempty"`
	Config   *Config `json:"config,omitempty"`
}

type Config struct {
	Scale int `json:"scale,omitempty"`
	//Ports   []string `json:"ports,omitempty"`
	//Env     []string `json:"env,omitempty"`
	//Volumes []string `json:"volumes,omitempty"`
}

func DeployCmd(name, image, template, url string, scale int) {

	var (
		ctx    = context.Get()
		config *Config
	)

	if scale != 0 /* || len(env) != 0 || len(ports) != 0 || len(volumes) != 0 */ {
		config = new(Config)
		config.Scale = scale
		//config.Env = env
		//config.Ports = ports
		//config.Volumes = volumes
	}

	err := Deploy(name, image, template, url, config)
	if err != nil {
		ctx.Log.Error(err)
		return
	}
}

func Deploy(name, image, template, url string, config *Config) error {

	var (
		err     error
		ctx     = context.Get()
		project = new(model.Project)
		er      = new(e.Http)
		res     = new(struct{})
	)

	err = ctx.Storage.Get("project", project)
	if err != nil {
		return errors.New(err.Error())
	}

	if project.ID == "" {
		return errors.New("Project didn't select")
	}

	var cfg = deployS{}
	cfg.Project = project.ID

	if name != "" {
		cfg.Name = name
	}

	if template != "" {
		cfg.Template = template
	}

	if image != "" {
		cfg.Image = image
	}

	if url != "" {
		cfg.Url = url
	}

	if config != nil {
		cfg.Config = config
	}

	_, _, err = ctx.HTTP.
		POST("/deploy").
		AddHeader("Content-Type", "application/json").
		AddHeader("Authorization", "Bearer "+ctx.Token).
		BodyJSON(cfg).
		Request(res, er)
	if err != nil {
		return err
	}

	if er.Code != 0 {
		return errors.New(e.Message(er.Status))
	}

	return nil
}
