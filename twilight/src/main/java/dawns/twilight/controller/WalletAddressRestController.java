package dawns.twilight.controller;

import java.util.List;

import javax.servlet.http.HttpServletRequest;

import org.apache.shiro.authz.annotation.RequiresAuthentication;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RequestParam;

import com.alibaba.fastjson.JSONObject;

import dawns.twilight.common.base.BaseRestController;
import dawns.twilight.common.base.Constants;
import dawns.twilight.common.chain.FabricService;
import dawns.twilight.common.web.JsonResult;
import dawns.twilight.common.web.RequestWallet;
import dawns.twilight.common.web.ResponseWallet;
import dawns.twilight.dao.model.WalletAddress;
import dawns.twilight.dao.model.WalletAddressExample;
import dawns.twilight.service.WalletAddressService;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import lombok.extern.slf4j.Slf4j;

@Controller
@Slf4j
@RequestMapping("/wallet")
@Api(value = "wallet api", tags = {"2. Wallet"})
public class WalletAddressRestController extends BaseRestController{
    @Autowired
    private WalletAddressService walletAddressService;

    @Autowired
    private FabricService fabric;
    
    @ApiOperation(value="register wallet")
    @RequestMapping(value = "", method = RequestMethod.POST)
    @RequiresAuthentication
    public JsonResult<ResponseWallet> register(HttpServletRequest request, @RequestBody RequestWallet req) {
    	Integer userId = (Integer) request.getAttribute(Constants.CURRENT_USER_ID);
    	log.debug("call register...");
    	JsonResult<ResponseWallet> result = new JsonResult<>(HttpStatus.OK);
    	String ret = fabric.NewAccount(String.valueOf(userId), req.getPass(),
    			req.getNetwork(), req.getToken());
    	JSONObject obj = JSONObject.parseObject(ret);
    	result.setCode(obj.getInteger("code"));
    	result.setMessage(obj.getString("message"));
    	ResponseWallet wallet = new ResponseWallet();
    	wallet.setNetwork(req.getNetwork());
    	wallet.setToken(req.getToken());
    	wallet.setHeight(obj.getInteger("height"));
    	wallet.setTxid(obj.getString("txid"));
    	wallet.setAddress(obj.getString("address"));
    	result.setData(wallet);
        return result;
    }
    
    @ApiOperation(value="query wallet")
    @RequestMapping(value = "/{id}", method = RequestMethod.GET)
    @RequiresAuthentication
    public JsonResult<WalletAddress> get(HttpServletRequest request, @PathVariable("id") Integer id) {
    	String chainRet = fabric.ChaincodeQuery("orgchannel", "wallet",
				"{\"Func\":\"query\", \"Args\":[\"sequence\"]}");
    	JsonResult<WalletAddress> ret = new JsonResult<>(HttpStatus.NOT_FOUND);
    	ret.setMessage(chainRet);
    	return ret;
    }

    @ApiOperation(value="query transaction")
    @RequestMapping(value = "/transaction/{sequence}", method = RequestMethod.GET)
    public JsonResult<WalletAddress> getTransaction(HttpServletRequest request,
    		@PathVariable("sequence") String sequence) {
    	String chainRet = fabric.ChaincodeQuery("orgchannel", "wallet",
				"{\"Func\":\"query\", \"Args\":[\"transaction\", \"sequence\", \"+sequence+\"]}");
    	JsonResult<WalletAddress> ret = new JsonResult<>(HttpStatus.NOT_FOUND);
    	ret.setMessage(chainRet);
    	return ret;
    }

    @ApiOperation(value="paging query")
    @RequestMapping(value = "", method = RequestMethod.GET)
    public JsonResult<List<WalletAddress>> page(HttpServletRequest request,
                                            @RequestParam(value = "pageNum") Integer pageNum,
                                            @RequestParam(value = "pageSize") Integer pageSize) {
        WalletAddressExample walletAddressExample=new WalletAddressExample();
        walletAddressExample.setOrderByClause("id");
        return new JsonResult<>(walletAddressService.selectByExampleForStartPage(walletAddressExample, pageNum,pageSize));
    }
}
