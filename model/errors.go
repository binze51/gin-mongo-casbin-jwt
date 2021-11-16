package model

import (
	"gitee.com/binze/binzekeji/pkg/response"
)

//账号体系错误码
const (
	RepeatedAccountName response.RetCode = 2001
	AccountCreateError  response.RetCode = 2002
	AccountUpdateError  response.RetCode = 2003
	AccountSearchError  response.RetCode = 2004

	AccountLoginError                response.RetCode = 2005
	AccountLoginLimitError           response.RetCode = 2006
	AccountLoginDenyError            response.RetCode = 2007
	IssueTokenFaildError             response.RetCode = 2008
	AccountPlatformLoginRepeatError  response.RetCode = 2009
	AccountLogoutError               response.RetCode = 2010
	AccountErrorLoginSaveError       response.RetCode = 2011
	AccountUpdateSaveError           response.RetCode = 2012
	AccountChangePwdError            response.RetCode = 2013
	AccountNewPwdSetFailError        response.RetCode = 2014
	OtherRepeatTokenMoveError        response.RetCode = 2015
	OtherRepeatRefreshTokeClearError response.RetCode = 2016
	RefreshTokenFaildError           response.RetCode = 2017
	RefreshTokenExpiredError         response.RetCode = 2018
	AdminUnlockError                 response.RetCode = 2019
	CorporationPostError             response.RetCode = 2020
	LoginMoveDenyError               response.RetCode = 2021
	UpdateEmailError                 response.RetCode = 2022
	UpdateMobileError                response.RetCode = 2023
	UpdateAvatarError                response.RetCode = 2024
	NoAdminUpdateAdminError          response.RetCode = 2025
	AccountErrorLoginCountError      response.RetCode = 2026
	ClearErrorLoginCountError        response.RetCode = 2027
)

var AccountCodeTable = map[response.RetCode]string{
	RepeatedAccountName: "账号名称重复",
	AccountCreateError:  "账号创建失败",
	AccountUpdateError:  "账号更新失败",
	AccountSearchError:  "账户查询失败",

	AccountLoginError:      "账号名或密码不正确，请重新输入!",
	AccountLoginLimitError: "账号名或密码不正确，还剩余%d次输入机会!",
	AccountLoginDenyError:  "输入错误次数已达上限，账号已锁定，请联系管理员解锁!",

	IssueTokenFaildError:             "签发授权token失败",
	AccountPlatformLoginRepeatError:  "该账号在其设备已登录,是否强制踢掉并登录",
	AccountLogoutError:               "该账号退出失败",
	AccountErrorLoginSaveError:       "首次登录错误记录失败",
	AccountUpdateSaveError:           "账号更新写库失败",
	AccountChangePwdError:            "更新密码旧密码错误",
	AccountNewPwdSetFailError:        "密码更新失败",
	OtherRepeatTokenMoveError:        "强制失效token失败",
	OtherRepeatRefreshTokeClearError: "清除所有token信息失败",
	RefreshTokenFaildError:           "刷新token失败",
	RefreshTokenExpiredError:         "刷新token已过期",
	AdminUnlockError:                 "管理员解锁失败",
	CorporationPostError:             "客户企业信息补录失败",
	LoginMoveDenyError:               "账号锁定失败",
	UpdateEmailError:                 "更新邮箱失败",
	UpdateMobileError:                "更新手机失败",
	UpdateAvatarError:                "更新头像失败",
	NoAdminUpdateAdminError:          "非admin账号不能更新admin账号",
	AccountErrorLoginCountError:      "账号登录错误计数失败",
	ClearErrorLoginCountError:        "清除账号错误登录计数失败",
}

func init() {
	response.MergeCodeTable(AccountCodeTable)
}
