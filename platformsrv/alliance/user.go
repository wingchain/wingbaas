
package alliance

import (
	"fmt"
	"time"
	"encoding/json"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/utils"
)

type User struct {
	Mail			string  	`json:"Mail"`
	Password		string  	`json:"Password"`
	VerifyCode		string  	`json:"VerifyCode,omitempty"` 
	Alliances		[]Alliance  `json:"Alliances,omitempty"`
}

const (
	USER_FILE	string = "users.json"
)

func AddUser(user User) error { 
	cfgPath := utils.BAAS_CFG.AlliancePath + USER_FILE
	exsist,_ := utils.PathExists(cfgPath) 
	var users []User
	if exsist {
		bytes,err := utils.LoadFile(cfgPath)
		if err == nil {
			err = json.Unmarshal(bytes,&users)
			if err != nil {
				logger.Errorf("AddUser: unmarshal users error,%v", err)
				return fmt.Errorf("%v", err)
			}
		}else {
			logger.Errorf("AddUser: load user list file error,%v", err)
			return fmt.Errorf("%v", err)
		}
	}
	for _,u := range users {
		if u.Mail == user.Mail {
			logger.Errorf("%s%s","AddUser: user mail already exsist ",user.Mail)
			return fmt.Errorf("%s%s","AddUser: user mail already exsist ",user.Mail)
		}
	}
	users = append(users,user)
	bytes, err := json.Marshal(users)
	if err != nil {
		logger.Errorf("AddUser: marshal users error,%v", err)
		return fmt.Errorf("%v", err)
	}
	err = utils.WriteFile(cfgPath,string(bytes))
	if err != nil {
		logger.Errorf("AddUser: Write user list file error,%v", err)
		return fmt.Errorf("%v", err)
	}
	return nil
}

func GetUsersByMail(mail string)(*User,error) {
	cfgPath := utils.BAAS_CFG.AlliancePath + USER_FILE
	exsist,_ := utils.PathExists(cfgPath)
	if exsist {
		var users []User
		bytes,err := utils.LoadFile(cfgPath)
		if err == nil {
			err = json.Unmarshal(bytes,&users)
			if err != nil {
				logger.Errorf("GetUsersByMail: unmarshal users error,%v", err)
				return nil,fmt.Errorf("%v", err)
			}
			for _,u := range users {
				if u.Mail == mail {
					return &u,nil
				}
			}
			logger.Errorf("GetUsersByMail: user not exsist")
			return nil,fmt.Errorf("GetUsersByMail: user not exsist")
		}else {
			logger.Errorf("GetUsersByMail: load user list file error,%v", err)
			return nil,fmt.Errorf("%v", err)
		}
	}  
	logger.Debug("GetUsersByMail: not find user list file")
	return nil,fmt.Errorf("GetUsersByMail: not find user list file")
}
 
func GetUsersByAllianceId(allianceId string)([]User,error) {
	cfgPath := utils.BAAS_CFG.AlliancePath + USER_FILE
	exsist,_ := utils.PathExists(cfgPath)
	if exsist {
		var users []User
		bytes,err := utils.LoadFile(cfgPath)
		if err == nil {
			err = json.Unmarshal(bytes,&users)
			if err != nil {
				logger.Errorf("GetUsersByAlliance: unmarshal users error,%v", err)
				return nil,fmt.Errorf("%v", err)
			}
			var us []User
			for _,u := range users {
				tmpUser := u
				for _,a := range u.Alliances {
					if a.Id == allianceId {
						us = append(us,tmpUser)
						break
					}
				}
			}
			return us,nil
		}else {
			logger.Errorf("GetUsersByAlliance: load user list file error,%v", err)
			return nil,fmt.Errorf("%v", err)
		}
	}
	logger.Debug("GetUsersByAlliance: not find user list file")
	return nil,fmt.Errorf("GetUsersByAlliance: not find user list file")
} 

func GetUserByMailAndPw(mail string,password string) (*User,error) { 
	cfgPath := utils.BAAS_CFG.AlliancePath + USER_FILE
	exsist,_ := utils.PathExists(cfgPath)
	if exsist {
		var users []User
		bytes,err := utils.LoadFile(cfgPath)
		if err == nil {
			err = json.Unmarshal(bytes,&users)
			if err != nil {
				logger.Errorf("GetUserByMailAndPw: unmarshal users error,%v", err)
				return nil,fmt.Errorf("%v", err)
			}
			for _,u := range users {
				if u.Mail == mail && u.Password == password {
					return &u,nil
				}
			}
			return nil,fmt.Errorf("GetUserByMailAndPw: user match failed")
		}else {
			logger.Errorf("GetUserByMailAndPw: load user list file error,%v", err)
			return nil,fmt.Errorf("%v", err)
		}
	}
	logger.Debug("GetUserByMailAndPw: not find user list file")
	return nil,fmt.Errorf("GetUserByMailAndPw: not find user list file")
}

func UserAddAlliance(mail string,alliance Alliance)error {
	_,err := GetAllianceById(alliance.Id)
	if err != nil {
		logger.Errorf("UserAddAlliance: not exsist the alliance,id=" + alliance.Id)
		return fmt.Errorf("UserAddAlliance: not exsist the alliance,id=" + alliance.Id)
	}
	cfgPath := utils.BAAS_CFG.AlliancePath + USER_FILE
	exsist,_ := utils.PathExists(cfgPath)
	if exsist {
		var users []User
		bytes,err := utils.LoadFile(cfgPath)
		if err == nil {
			err = json.Unmarshal(bytes,&users)
			if err != nil {
				logger.Errorf("UserAddAlliance: unmarshal users error,%v", err)
				return fmt.Errorf("%v", err)
			}
			var us []User
			find := false
			for _,u := range users {
				tmpUser := u
				if u.Mail == mail {
					for _,a := range u.Alliances {
						if a.Id == alliance.Id {
							logger.Errorf("UserAddAlliance: user already in this alliance")
							return fmt.Errorf("UserAddAlliance: user already in this alliance")
						}
					}
					nowTime := time.Now().Format("2006-01-02 15:04:05")
					alliance.JoinTime = nowTime
					tmpUser.Alliances = append(tmpUser.Alliances,alliance)
					find = true
				}
				us = append(us,tmpUser)
			}
			if !find {
				logger.Errorf("UserAddAlliance: user not exsist")
				return fmt.Errorf("UserAddAlliance: user not exsist")
			}
			bytes, err := json.Marshal(us)
			if err != nil {
				logger.Errorf("UserAddAlliance: marshal users error,%v", err)
				return fmt.Errorf("%v", err)
			}
			err = utils.WriteFile(cfgPath,string(bytes))
			if err != nil {
				logger.Errorf("UserAddAlliance: Write users file error,%v", err)
				return fmt.Errorf("%v", err)
			}
			return nil
		}else {
			logger.Errorf("UserAddAlliance: load user list file error,%v", err)
			return fmt.Errorf("%v", err)
		}
	}
	logger.Errorf("UserAddAlliance: not find user list file")
	return fmt.Errorf("UserAddAlliance: not find user list file")
}

func GetUserCreatedAlliances(mail string)([]Alliance,error) {
	cfgPath := utils.BAAS_CFG.AlliancePath + USER_FILE
	exsist,_ := utils.PathExists(cfgPath)
	if exsist {
		var users []User
		bytes,err := utils.LoadFile(cfgPath)
		if err == nil {
			err = json.Unmarshal(bytes,&users)
			if err != nil {
				logger.Errorf("GetUserCreatedAlliances: unmarshal users error,%v", err)
				return nil,fmt.Errorf("%v", err)
			}
			var as []Alliance
			for _,u := range users {
				if u.Mail == mail {
					for _,a := range u.Alliances {
						tmpAlliance := a
						if a.Creator == mail {
							as = append(as,tmpAlliance)
						}
					}
					
				}
			}
			return as,nil
		}else {
			logger.Errorf("GetUserCreatedAlliances: load user list file error,%v", err)
			return nil,fmt.Errorf("%v", err)
		}
	}  
	logger.Errorf("GetUserCreatedAlliances: not find user list file")
	return nil,fmt.Errorf("GetUserCreatedAlliances: not find user list file")
}

func GetUserJoinedAlliances(mail string)([]Alliance,error) {
	cfgPath := utils.BAAS_CFG.AlliancePath + USER_FILE
	exsist,_ := utils.PathExists(cfgPath)
	if exsist {
		var users []User
		bytes,err := utils.LoadFile(cfgPath)
		if err == nil {
			err = json.Unmarshal(bytes,&users)
			if err != nil {
				logger.Errorf("GetUserJoinedAlliances: unmarshal users error,%v", err)
				return nil,fmt.Errorf("%v", err)
			}
			for _,u := range users { 
				if u.Mail == mail { 
					return u.Alliances,nil
				}
			}
			logger.Errorf("GetUserJoinedAlliances: not find user")
			return nil,fmt.Errorf("GetUserJoinedAlliances: not find user")
		}else {
			logger.Errorf("GetUserJoinedAlliances: load user list file error,%v", err)
			return nil,fmt.Errorf("%v", err)
		}
	}  
	logger.Errorf("GetUserJoinedAlliances: not find user list file")
	return nil,fmt.Errorf("GetUserJoinedAlliances: not find user list file")
}


func DeleteUserAlliance(allianceId string) error {
	cfgPath := utils.BAAS_CFG.AlliancePath + USER_FILE
	exsist,_ := utils.PathExists(cfgPath)
	if exsist {
		var users []User
		var us []User
		bytes,err := utils.LoadFile(cfgPath)
		if err == nil {
			err = json.Unmarshal(bytes,&users)
			if err != nil {
				logger.Errorf("DeleteUserAlliance: unmarshal users error,%v", err)
				return fmt.Errorf("%v", err)
			}
			for _,u := range users {
				var as []Alliance
				tmpUser := u
				for _,a := range u.Alliances {
					tmpAlliance := a
					if tmpAlliance.Id != allianceId {	
						as = append(as,tmpAlliance)
					}
				}
				tmpUser.Alliances = as
				us = append(us,tmpUser)
			}  
			bytes, err := json.Marshal(us)
			if err != nil {
				logger.Errorf("DeleteUserAlliance: marshal users error,%v", err)
				return fmt.Errorf("%v", err)
			}
			err = utils.WriteFile(cfgPath,string(bytes))
			if err != nil {  
				logger.Errorf("DeleteUserAlliance: Write users file error,%v", err)
				return fmt.Errorf("%v", err)
			} 
			return nil 
		}else {
			logger.Errorf("DeleteUserAlliance: load user list file error,%v", err)
			return fmt.Errorf("%v", err)
		}
	}  
	logger.Errorf("DeleteUserAlliance: not find user list file")
	return fmt.Errorf("DeleteUserAlliance: not find user list file")
}

func DeleteUserJoinedAlliance(mail string,allianceId string) error {
	cfgPath := utils.BAAS_CFG.AlliancePath + USER_FILE
	exsist,_ := utils.PathExists(cfgPath)
	if exsist {
		var users []User
		var us []User
		bytes,err := utils.LoadFile(cfgPath)
		if err == nil {
			err = json.Unmarshal(bytes,&users)
			if err != nil {
				logger.Errorf("DeleteUserJoinedAlliance: unmarshal users error,%v", err)
				return fmt.Errorf("%v", err)
			}
			var as []Alliance
			for _,u := range users {
				tmpUser := u
				if u.Mail != mail {
					us = append(us,tmpUser)
					continue
				}else {
					for _,a := range u.Alliances {
						tmpAlliance := a
						if a.Id != allianceId {
							as = append(as,tmpAlliance)
						}
					}
					tmpUser.Alliances = as
					us = append(us,tmpUser)
				}	
			} 
			bytes, err := json.Marshal(us)
			if err != nil {
				logger.Errorf("DeleteUserJoinedAlliance: marshal users error,%v", err)
				return fmt.Errorf("%v", err)
			}
			err = utils.WriteFile(cfgPath,string(bytes))
			if err != nil {  
				logger.Errorf("DeleteUserJoinedAlliance: Write users file error,%v", err)
				return fmt.Errorf("%v", err)
			} 
			return nil 
		}else {
			logger.Errorf("DeleteUserJoinedAlliance: load user list file error,%v", err)
			return fmt.Errorf("%v", err)
		}
	}  
	logger.Errorf("DeleteUserJoinedAlliance: not find user list file")
	return fmt.Errorf("DeleteUserJoinedAlliance: not find user list file")
}