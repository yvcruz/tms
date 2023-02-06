package tms

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type TodusMessageService struct {
	Config  TodusMessageServiceConfig
	Token   string
	refresh string
}

type TodusMessageServiceConfig struct {
	Url      string
	Username string
	Password string
	Uid      string
}

func NewTodusMessageService(config TodusMessageServiceConfig) *TodusMessageService {
	return &TodusMessageService{
		Config:  config,
		Token:   "",
		refresh: "",
	}
}

func (tms TodusMessageService) GetToken() string {
	token, refresh, err := tms.getTokenFromFile()
	if err != nil {
		if tk, err := tms.getTokenFromUrl(); err != nil {
			return ""
		} else {
			return tk
		}
	} else {
		if len(token) == 0 || len(refresh) == 0 {
			if tk, err := tms.getTokenFromUrl(); err != nil {
				return ""
			} else {
				return tk
			}
		} else {
			tms.Token = token
			tms.refresh = refresh
			if tk, err := tms.validateToken(); err != nil {
				return ""
			} else {
				return tk
			}
		}
	}
}

func (tms TodusMessageService) validateToken() (string, error) {
	url := fmt.Sprintf("%s/user/me", tms.Config.Url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", tms.Token)

	client := &http.Client{}
	r, err := client.Do(req)
	if err != nil {
		panic(err)
		return "", err
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		if err := tms.clearTokenFile(); err != nil {
			return "", err
		}
		if token, err := tms.refreshToken(); err != nil {
			return "", err
		} else {
			return token, nil
		}
	}
	return tms.Token, nil
}

func (tms TodusMessageService) getTokenFromUrl() (string, error) {
	postBody, _ := json.Marshal(map[string]string{
		"username": tms.Config.Username,
		"password": tms.Config.Password,
	})
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(tms.Config.Url+"/auth", "application/json", responseBody)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("status: %d", resp.StatusCode))
	}
	buf := make([]byte, 1024)
	num, err := resp.Body.Read(buf)
	res := string(buf[0:num])
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(res), &m); err != nil {
		return "", err
	}
	token := m["tk"].(string)
	refresh := m["rtk"].(string)

	if err = ioutil.WriteFile("./tk", []byte(token), 0644); err != nil {
		panic(err)
	}

	if err = ioutil.WriteFile("./rtk", []byte(refresh), 0644); err != nil {
		panic(err)
	}
	tms.Token = token
	tms.refresh = refresh
	return token, nil
}

func (tms TodusMessageService) refreshToken() (string, error) {
	if tms.refresh != "" {
		postBody, _ := json.Marshal(map[string]string{
			"rtk": tms.refresh,
		})
		responseBody := bytes.NewBuffer(postBody)

		resp, err := http.Post(tms.Config.Url+"/auth/refresh", "application/json", responseBody)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			buf := make([]byte, 1024)
			num, _ := resp.Body.Read(buf)
			res := string(buf[0:num])
			var m map[string]interface{}
			if err := json.Unmarshal([]byte(res), &m); err != nil {
				return "", err
			}
			token := m["tk"].(string)
			refresh := m["rtk"].(string)
			ftk, err := os.Create("./rtk")
			if err != nil {
				return "", err

			}
			defer ftk.Close()
			frtk, err := os.Create("./rtk")
			if err != nil {
				return "", err
			}
			defer frtk.Close()

			_, err = ftk.Write([]byte(token))
			if err != nil {
				return "", err
			}
			_, err = frtk.Write([]byte(refresh))
			if err != nil {
				return "", err
			}

			tms.refresh = refresh
			tms.Token = token
			return token, nil
		} else if resp.StatusCode == 401 {
			return tms.getTokenFromUrl()
		}
	} else {
		return tms.getTokenFromUrl()
	}
	return tms.getTokenFromUrl()
}

func (tms TodusMessageService) getTokenFromFile() (string, string, error) {
	tk, err := ioutil.ReadFile("./tk")
	if err != nil {
		return "", "", err
	}
	rtk, err := ioutil.ReadFile("./rtk")
	if err != nil {
		return "", "", err
	}
	return string(tk), string(rtk), nil
}

func (tms TodusMessageService) clearTokenFile() error {
	err := os.Remove("./tk")
	if err != nil {
		return err
	}
	err = os.Remove("./rtk")
	if err != nil {
		return err
	}
	return nil
}

func (tms TodusMessageService) SendMessageToGroup(message string) bool {
	uid := tms.Config.Uid
	url := fmt.Sprintf("%s/sendgroup", tms.Config.Url)

	tk := tms.GetToken()
	if len(tk) == 0 {
		return false
	}

	postBody, _ := json.Marshal(map[string]interface{}{
		"to":   uid,
		"body": message,
		"from": tms.Config.Username,
		"type": 1,
	})

	posBodyBuffer := bytes.NewBuffer(postBody)
	req, err := http.NewRequest("PUT", url, posBodyBuffer)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tk)
	client := &http.Client{}
	r, err := client.Do(req)
	if err != nil {
		panic(err)
		return false
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		return false
	}
	return false
}

func (tms TodusMessageService) SendMessageToUser(username, message string) bool {
	url := fmt.Sprintf("%s/msg", tms.Config.Url)

	tk := tms.GetToken()
	if len(tk) == 0 {
		return false
	}

	postBody, _ := json.Marshal(map[string]interface{}{
		"to":   username,
		"body": message,
		"from": tms.Config.Username,
	})

	posBodyBuffer := bytes.NewBuffer(postBody)
	req, err := http.NewRequest("PUT", url, posBodyBuffer)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tk)
	client := &http.Client{}
	r, err := client.Do(req)
	if err != nil {
		panic(err)
		return false
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		if r.StatusCode == 405 {
			log.Println("err: user not allowed")
		}
		return false
	}
	return false
}
