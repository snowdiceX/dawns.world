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
import dawns.twilight.common.cazimi.CazimiService;
import dawns.twilight.common.web.JsonResult;
import dawns.twilight.common.web.RequestWallet;
import dawns.twilight.dao.model.WalletAddress;
import dawns.twilight.dao.model.WalletAddressExample;
import dawns.twilight.service.WalletAddressService;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import lombok.extern.slf4j.Slf4j;

@Controller
@Slf4j
@RequestMapping("/wallet")
@Api(value = "WalletAddress api", tags = {"2. Wallet"})
public class WalletAddressRestController extends BaseRestController{
    @Autowired
    private WalletAddressService walletAddressService;

    @Autowired
    private CazimiService cazimi;
    
    @ApiOperation(value="请求创建托管钱包地址")
    @RequestMapping(value = "", method = RequestMethod.POST)
    @RequiresAuthentication
    public JsonResult<String> register(HttpServletRequest request, @RequestBody RequestWallet req) {
    	Integer userId = (Integer) request.getAttribute(Constants.CURRENT_USER_ID);
    	log.debug("call register...");
    	JsonResult<String> result = new JsonResult<>(HttpStatus.OK);
    	String txid = "txid";
    	String ret = cazimi.registerWallet(String.valueOf(userId), req.getToken(), req.getNetwork());
    	JSONObject obj = JSONObject.parseObject(ret);
    	result.setCode(obj.getInteger("code"));
    	result.setMessage(obj.getString("message"));
    	result.setData(txid);
        return result;
    }

    @ApiOperation(value="根据id更新WalletAddress")
    @RequestMapping(value = "", method = RequestMethod.PUT)
    public JsonResult<Integer> update(HttpServletRequest request, @RequestBody WalletAddress walletAddress) {
        return new JsonResult<>(walletAddressService.updateByPrimaryKey(walletAddress));
    }

    @ApiOperation(value="根据id查询WalletAddress")
    @RequestMapping(value = "/{id}", method = RequestMethod.GET)
    public JsonResult<WalletAddress> get(HttpServletRequest request, @PathVariable("id") Integer id) {
        WalletAddress walletAddress=walletAddressService.selectByPrimaryKey(id);
        if(walletAddress!=null){
            return new JsonResult<>(walletAddress);
        }else{
            return new JsonResult<>(HttpStatus.NOT_FOUND);
        }
    }

    @ApiOperation(value="根据id删除WalletAddress")
    @RequestMapping(value = "/{id}", method = RequestMethod.DELETE)
    public JsonResult<Integer> delete(HttpServletRequest request, @PathVariable("id") Integer id) {
        return new JsonResult<>(walletAddressService.deleteByPrimaryKey(id));
    }

    @ApiOperation(value="分页查询WalletAddress")
    @RequestMapping(value = "", method = RequestMethod.GET)
    public JsonResult<List<WalletAddress>> page(HttpServletRequest request,
                                            @RequestParam(value = "pageNum") Integer pageNum,
                                            @RequestParam(value = "pageSize") Integer pageSize) {
        WalletAddressExample walletAddressExample=new WalletAddressExample();
        walletAddressExample.setOrderByClause("id");
        return new JsonResult<>(walletAddressService.selectByExampleForStartPage(walletAddressExample, pageNum,pageSize));
    }
}
