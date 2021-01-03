package services

import (
	"log"
	"strings"
	"tbot/config"
	"tbot/internal/model"
)

type ServiceManager struct {
	cfg      *config.Config
	Services map[string]Service
}

func NewServiceManager(cfg *config.Config, services ...Service) *ServiceManager {
	serviceManager := ServiceManager{
		cfg:      cfg,
		Services: make(map[string]Service),
	}
	for _, service := range services {
		serviceManager.Services[service.GetTag()] = service
	}
	return &serviceManager
}

func (sm *ServiceManager) Act(msgText string) (model.Resp, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("recovered from ", r)
		}
	}()

	if len(msgText) < 3 {
		return model.Resp{Text: msgText}, nil
	}

	for serviceTag, acts := range sm.cfg.Stg.Services {
		for _, act := range acts {
			if strings.HasPrefix(msgText, act) {
				text := msgText[len(act):]
				return sm.chooseAction(serviceTag, text)
			}
		}
	}

	if strings.HasPrefix(msgText, "/help") {
		return model.Resp{Text: sm.cfg.Stg.HelpText}, nil
	}
	return model.Resp{Text: msgText}, nil
}

func (sm *ServiceManager) chooseAction(serviceTag string, text string) (model.Resp, error) {
	if text == "" {
		return model.Resp{Text: text}, nil
	}
	switch serviceTag {
	case sm.cfg.Stg.WikiStg.Tag:
		return sm.takeAction(sm.cfg.Stg.WikiStg.Tag, text)
	case sm.cfg.Stg.NewtonStg.Tag:
		return sm.takeAction(sm.cfg.Stg.NewtonStg.Tag, text)
	case sm.cfg.Stg.PlaygroundStg.Tag:
		return sm.takeAction(sm.cfg.Stg.PlaygroundStg.Tag, text)
	case sm.cfg.Stg.MarkovStg.Tag:
		return sm.takeAction(sm.cfg.Stg.MarkovStg.Tag, text)
	case sm.cfg.Stg.GolangStg.Tag:
		return sm.takeAction(sm.cfg.Stg.GolangStg.Tag, text)
	}
	return model.Resp{Text: text}, nil
}

func (sm *ServiceManager) takeAction(tag string, text string) (model.Resp, error) {
	if sm.Services[tag].IsReady() {
		result, err := sm.Services[tag].Query(text)
		if err != nil {
			return model.Resp{}, err
		}
		return result, nil
	}
	return model.Resp{Text: "Not ready yet. Please, waiting some minutes"}, nil
}
