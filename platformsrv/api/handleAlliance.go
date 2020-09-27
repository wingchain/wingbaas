
package api

import (
	"encoding/json"
	"net/http"
	"io/ioutil"
	"github.com/labstack/echo/v4"
	"github.com/wingbaas/platformsrv/logger"
	"github.com/wingbaas/platformsrv/alliance"
)

func userRegister(c echo.Context) error {
	logger.Debug("userRegister")
	var user alliance.User
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	result, err := ioutil.ReadAll(c.Request().Body) 
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
    }
    err = json.Unmarshal(result, &user)
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	err = alliance.AddUser(user)
	if err != nil {
        msg := err.Error()
        ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
    }
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,nil)
	return c.JSON(http.StatusOK,ret) 
}

func userLogin(c echo.Context) error {
	logger.Debug("userLogin")
	var user alliance.User 
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	result, err := ioutil.ReadAll(c.Request().Body) 
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
    }
    err = json.Unmarshal(result, &user)
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	_,err = alliance.GetUserByMailAndPw(user.Mail,user.Password)
	if err != nil {
        msg := err.Error()
        ret := getApiRet(CODE_ERROR_EXE,msg,nil) 
		return c.JSON(http.StatusOK,ret)
    }
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,nil)
	return c.JSON(http.StatusOK,ret) 
}

func creatAlliance(c echo.Context) error {
	logger.Debug("creatAlliance")
	var a alliance.Alliance
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	result, err := ioutil.ReadAll(c.Request().Body) 
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
    }
    err = json.Unmarshal(result, &a)
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	_,err = alliance.GetUsersByMail(a.Creator) 
	if err != nil {
        // msg := err.Error()
        // ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		// return c.JSON(http.StatusOK,ret) 

		//user not exsist,new user
		var user alliance.User
		user.Mail = a.Creator 
		err = alliance.AddUser(user)
		if err != nil {
       		msg := err.Error()
        	ret := getApiRet(CODE_ERROR_EXE,msg,nil)
			return c.JSON(http.StatusOK,ret)
    	} 
	}
	id,err := alliance.AppendAlliance(a) 
	if err != nil {
        msg := err.Error()
        ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	a.Id = id
	err = alliance.UserAddAlliance(a.Creator,a)
	if err != nil {
        msg := err.Error()
        ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret) 
	}
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,id)
	return c.JSON(http.StatusOK,ret) 
}

func getAlliance(c echo.Context) error {
	logger.Debug("getAlliance")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	as,err := alliance.GetAlliances()
	if err != nil {
        msg := err.Error()
        ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
    }
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,as)
	return c.JSON(http.StatusOK,ret) 
}

func getAllianceById(c echo.Context) error {
	logger.Debug("getAllianceById")
	allianceId := c.Param("allianceid")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	a,err := alliance.GetAllianceById(allianceId)
	if err != nil {
        msg := err.Error()
        ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
    }
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,*a)
	return c.JSON(http.StatusOK,ret) 
}

func userAddAlliance(c echo.Context) error {
	logger.Debug("userAddAlliance")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	type ParaSt struct {
		Mail string `json:"Mail"`
		Alliance alliance.Alliance `json:"Alliance"`
	}
	var obj ParaSt 
	result, err := ioutil.ReadAll(c.Request().Body) 
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
    }
    err = json.Unmarshal(result, &obj)
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	//user not exsist,new user
	_,err = alliance.GetUsersByMail(obj.Mail) 
	if err != nil {
		var user alliance.User
		user.Mail = obj.Mail 
		err = alliance.AddUser(user)
		if err != nil {
       		msg := err.Error()
        	ret := getApiRet(CODE_ERROR_EXE,msg,nil)
			return c.JSON(http.StatusOK,ret)
    	} 
	}
	err = alliance.UserAddAlliance(obj.Mail,obj.Alliance)
	if err != nil {
        msg := err.Error()
        ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
    }
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,nil)
	return c.JSON(http.StatusOK,ret) 
} 

func getJoinedAlliance(c echo.Context) error {
	logger.Debug("getJoinedAlliance")
	userMail := c.Param("usermail")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	as,err := alliance.GetUserJoinedAlliances(userMail)
	if err != nil {
        msg := err.Error()
        ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
    }
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,as)
	return c.JSON(http.StatusOK,ret) 
}

func getCreatedAlliance(c echo.Context) error {
	logger.Debug("getCreatedAlliance")
	userMail := c.Param("usermail")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	as,err := alliance.GetUserCreatedAlliances(userMail)
	if err != nil {
        msg := err.Error()
        ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
    }
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,as)
	return c.JSON(http.StatusOK,ret) 
}

func deleteAlliance(c echo.Context) error {
	logger.Debug("deleteAlliance")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	type ParaSt struct {
		Mail string `json:"Mail"`
		AllianceId string`json:"AllianceId"`
	}
	var obj ParaSt 
	result, err := ioutil.ReadAll(c.Request().Body) 
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
    }
    err = json.Unmarshal(result, &obj)
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	err = alliance.DeleteAllianceById(obj.Mail,obj.AllianceId)
	if err != nil {
        msg := err.Error()
        ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
    }
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,nil)
	return c.JSON(http.StatusOK,ret) 
} 

func getAllianceChains(c echo.Context) error {
	logger.Debug("getAllianceChains")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	allianceId := c.Param("allianceid")
	chs,err := alliance.GetAllianceChains(allianceId)
	if err != nil {
        msg := err.Error()
        ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
    } 
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,chs)
	return c.JSON(http.StatusOK,ret) 
} 

func getAllianceClusters(c echo.Context) error {
	logger.Debug("getAllianceClusters")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	allianceId := c.Param("allianceid")
	cs,err := alliance.GetAllianceClusters(allianceId)
	if err != nil {
        msg := err.Error()
        ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
    } 
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,cs)
	return c.JSON(http.StatusOK,ret) 
} 

func getUsersByAllianceId(c echo.Context) error {
	logger.Debug("getUsersByAllianceId")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	allianceId := c.Param("allianceid")
	chs,err := alliance.GetUsersByAllianceId(allianceId)
	if err != nil {
        msg := err.Error()
        ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
    } 
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,chs)
	return c.JSON(http.StatusOK,ret) 
}  

func deleteAllianceUser(c echo.Context) error {
	logger.Debug("deleteAllianceUser")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("content-type", "application/json")
	type ParaSt struct {
		Mail string `json:"Mail"`
		AllianceId string`json:"AllianceId"`
	}
	var obj ParaSt 
	result, err := ioutil.ReadAll(c.Request().Body) 
    if err != nil {
		msg := "read request body error"
		ret := getApiRet(CODE_ERROR_BODY,msg,nil)
		return c.JSON(http.StatusOK,ret)
    }
    err = json.Unmarshal(result, &obj)
    if err != nil {
        msg := "body json Unmarshal err"
        ret := getApiRet(CODE_ERROR_MASHAL,msg,nil)
		return c.JSON(http.StatusOK,ret)
	}
	err = alliance.DeleteUserJoinedAlliance(obj.Mail,obj.AllianceId)
	if err != nil {
        msg := err.Error()
        ret := getApiRet(CODE_ERROR_EXE,msg,nil)
		return c.JSON(http.StatusOK,ret)
    }
	ret := getApiRet(CODE_SUCCESS,MSG_SUCCESS,nil)
	return c.JSON(http.StatusOK,ret) 
} 
