package dawns.twilight.controller;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import org.apache.shiro.authz.annotation.RequiresAuthentication;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.ResponseBody;

import com.alibaba.fastjson.JSONObject;

import dawns.twilight.common.base.BaseRestController;
import dawns.twilight.common.base.Constants;
import dawns.twilight.common.chain.FabricService;
import dawns.twilight.common.web.JsonResult;
import dawns.twilight.common.web.RegisterTransaction;
import dawns.twilight.common.web.RequestWallet;
import dawns.twilight.common.web.ResponseWallet;
import dawns.twilight.dao.model.WalletAddress;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import lombok.extern.slf4j.Slf4j;

@Controller
@Slf4j
@RequestMapping("/wallet")
@Api(value = "wallet api", tags = {"2. Wallet"})
public class WalletAddressRestController extends BaseRestController{

    @Autowired
    private FabricService fabric;
    
    @ApiOperation(value="Register wallet",
    		notes="注册托管钱包",
    		tags="Register")
    @RequestMapping(value = "", method = RequestMethod.POST)
    @RequiresAuthentication
    public JsonResult<ResponseWallet> register(HttpServletRequest request, @RequestBody RequestWallet req) {
    	Integer userId = (Integer) request.getAttribute(Constants.CURRENT_USER_ID);
    	log.debug("call register...");
    	JsonResult<ResponseWallet> result = new JsonResult<>(HttpStatus.OK);
    	String ret = fabric.RegisterWallet(String.valueOf(userId), req.getPass(),
    			req.getNetwork(), req.getToken());
    	JSONObject obj = JSONObject.parseObject(ret);
    	result.setCode(obj.getInteger("code"));
    	result.setMessage(obj.getString("message"));
    	ResponseWallet wallet = new ResponseWallet();
    	obj = obj.getJSONObject("result");
    	wallet.setNetwork(req.getNetwork());
    	wallet.setToken(req.getToken());
    	wallet.setHeight(obj.getString("height"));
    	wallet.setTxid(obj.getString("txid"));
    	wallet.setAddress(obj.getString("address"));
    	result.setData(wallet);
        return result;
    }
    
    @ApiOperation(value="Withdraw", notes="托管钱包提取")
    @RequestMapping(value = "/withdraw",
    		method = RequestMethod.POST,
    		produces="application/json;charset=UTF-8")
    @ResponseBody
    @RequiresAuthentication
    public String withdraw(HttpServletRequest request, @RequestBody RequestWallet req) {
    	return "";
    }
    
    @ApiOperation(value="Import wallet",
    		notes="！！！仅用于测试！！！导入托管钱包",
    		tags="Register")
    @RequestMapping(value = "",
    		method = RequestMethod.PUT,
    		produces="application/json;charset=UTF-8")
    @ResponseBody
    @RequiresAuthentication
    public String importWallet(
    		HttpServletRequest request, @RequestBody RequestWallet req) {
    	Integer userId = (Integer) request.getAttribute(Constants.CURRENT_USER_ID);
    	log.debug("call register...");
    	return fabric.ChaincodeInvoke("orgchannel", "wallet",
    			"{\"Func\":\"register\", \"Args\":[\"wallet\", \""+userId+"\",\""
    					+req.getAddress()+"\",\""+req.getNetwork()+"\",\""
    					+req.getToken()+"\",\""+req.getHeight()+"\"]}");
    }
    
    @ApiOperation(value="Query wallet")
    @RequestMapping(value = "/{chain}/{token}/{address}",
    		method = RequestMethod.GET,
    		produces="application/json;charset=UTF-8")
    @ResponseBody
    @RequiresAuthentication
    public String getWallet(HttpServletRequest request,
    		@PathVariable("chain") String chain,
    		@PathVariable("token") String token,
    		@PathVariable("address") String address) {
    	return fabric.ChaincodeQuery("orgchannel", "wallet",
				"{\"Func\":\"query\", \"Args\":[\"wallet\", \""
						+chain+"\", \""+token+"\", \""+address+"\"]}");
    }
    
    @ApiOperation(value="Paging wallet")
    @RequestMapping(value = "",
    		method = RequestMethod.GET,
    		produces="application/json;charset=UTF-8")
    @ResponseBody
    public String pageWallet(
    		@RequestParam(value = "pageNum") Integer pageNum,
    		@RequestParam(value = "pageSize") Integer pageSize) {
    	return "";
    }
    
    @ApiOperation(value="Query sequence")
    @RequestMapping(value = "/sequence", method = RequestMethod.GET)
    @RequiresAuthentication
    public JsonResult<WalletAddress> getSequence(HttpServletRequest request) {
    	String chainRet = fabric.ChaincodeQuery("orgchannel", "wallet",
				"{\"Func\":\"query\", \"Args\":[\"sequence\"]}");
    	JsonResult<WalletAddress> ret = new JsonResult<>(HttpStatus.OK);
    	ret.setMessage(chainRet);
    	return ret;
    }

    @ApiOperation(value="Query fabric Tx")
    @RequestMapping(value = "/transaction/{sequence}", method = RequestMethod.GET)
    public JsonResult<WalletAddress> getTransaction(HttpServletRequest request,
    		@PathVariable("sequence") String sequence) {
    	String chainRet = fabric.ChaincodeQuery("orgchannel", "wallet",
				"{\"Func\":\"query\", \"Args\":[\"transaction\", \"sequence\", \"+sequence+\"]}");
    	JsonResult<WalletAddress> ret = new JsonResult<>(HttpStatus.OK);
    	ret.setMessage(chainRet);
    	return ret;
    }

    @ApiOperation(value="Paging registered Txs")
    @RequestMapping(value = "/transactions/{chain}/{token}/{walletAddress}",
    		method = RequestMethod.GET,
    		produces="application/json;charset=UTF-8")
    @ResponseBody
    public String pageTransactions(HttpServletRequest request, HttpServletResponse response,
    		@PathVariable("chain") String chain,
    		@PathVariable("token") String token,
    		@PathVariable("walletAddress") String walletAddress,
    		@RequestParam(value = "pageNum") Integer pageNum,
    		@RequestParam(value = "pageSize") Integer pageSize) {
    	String ret = fabric.ChaincodeQuery("orgchannel", "wallet",
				"{\"Func\":\"query\", \"Args\":[\"transaction\", \"page\","
				+ " \"registered\", \""+chain+"\", \""+token
				+"\", \"0x"+Integer.toHexString(pageNum)
				+"\", \"0x"+Integer.toHexString(pageSize)+"\", \""+walletAddress+"\"]}");
        return ret;
    }
    
    @ApiOperation(value="Register Tx",
    		notes="注册经过共识的外部区块链网络交易(非应用接口，由拥有授权证书的中继对网络和交易适配并共识后调用)",
    		tags="Register")
    @RequestMapping(value = "/transaction", method = RequestMethod.POST)
    public JsonResult<String> registerTransaction(
    		HttpServletRequest request, @RequestBody RegisterTransaction req) {
    	log.debug("call register transaction...");
    	JsonResult<String> result = new JsonResult<>(HttpStatus.OK);
    	String json = JSONObject.toJSONString(req);
    	json = json.replaceAll("\"", "\\\\\"");
    	log.info("register transaction: "+json);
    	String ret = fabric.ChaincodeInvoke("orgchannel", "wallet",
    			"{\"Func\":\"register\", \"Args\":[\"transaction\", \""+json+"\"]}");
    	JSONObject obj = JSONObject.parseObject(ret);
    	result.setCode(obj.getInteger("code"));
    	result.setMessage(obj.getString("message"));
    	result.setData(ret);
        return result;
    }
}
