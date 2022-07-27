/*******
* @Author:qingmeng
* @Description:
* @File:const
* @Date:2022/7/25
 */

package transfer

import "errors"

var (
	ERRPRIVILEGE=errors.New("拒绝连接")
	ERRDATA=errors.New("传输的数据错误")
	ERRSESSIONW=errors.New("会话写入错误")
	ERRSESSIONR=errors.New("会话读取错误")
	ERRENCODE=errors.New("编码错误")
	ERRDECODE=errors.New("解码错误")
	ERRTYPE =errors.New("类型错误")
	ERRNUM=errors.New("方法序号错误")
	ERRPARAMS=errors.New("方法参数错误")
	ERRTIME=errors.New("rpc连接超时")
)


