/**
 * @Description:
 * @Author: cnHuaShao
 * @File:  Msg
 * @Version: 1.0.0
 */

package model

// 消息模型
type Msg struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data interface{}
}

type Code struct {
	Code string
	Msg  string
}

var MsgOK = Code{
	Code: "0x000000",
	Msg:  "success",
}

/**
主消息生成函数
*/
func ReturnReqMess(code Code, data interface{}) *Msg {
	mess := Msg{
		Code: code.Code,
		Msg:  code.Msg,
		Data: data,
	}
	return &mess
}
