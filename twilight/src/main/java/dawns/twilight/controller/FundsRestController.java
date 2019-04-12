package dawns.twilight.controller;

import javax.servlet.http.HttpServletRequest;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.ResponseBody;

import dawns.twilight.common.base.BaseRestController;
import dawns.twilight.common.chain.FabricService;
import dawns.twilight.common.web.JsonResult;
import dawns.twilight.dao.model.Quotation;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;

@Controller
@RequestMapping("/funds")
@Api(value = "Funds api", description = "Funds api",tags = {"3. Funds"})
public class FundsRestController extends BaseRestController{
    
    @Autowired
    private FabricService fabric;

    @ApiOperation(value="Deposit",
    		notes="转入基金")
    @RequestMapping(value = "/deposit/{chain}/{network}/{address}", method = RequestMethod.POST)
    public JsonResult<Integer> add(HttpServletRequest request, @RequestBody Quotation quotation) {
        return new JsonResult<>(null);
    }

    @ApiOperation(value="Withdraw",
    		notes="转出基金")
    @RequestMapping(value = "/withdraw/{chain}/{network}/{address}", method = RequestMethod.POST)
    public JsonResult<Integer> delete(HttpServletRequest request, @PathVariable("id") Integer id) {
        return new JsonResult<>(null);
    }
    
    /**
     * @param request
     * @param chain
     * @param address
     * @return
     */
    @ApiOperation(value="Register funds",
    		notes="注册平台允许基金(通证对)，注册成功之后才可以通过唯一标识抵押和提取以赚取承兑服务费",
    		tags="Register")
    @RequestMapping(value = "/{chain}/{token}/base/{baseChain}/{baseToken}",
    		method = RequestMethod.POST,
    		produces="application/json;charset=UTF-8")
    @ResponseBody
    public String registerFunds(HttpServletRequest request,
    		@PathVariable("chain") String chain, @PathVariable("token") String token,
    		@PathVariable("baseChain") String baseChain, @PathVariable("baseToken") String baseToken) {
    	return fabric.RegisterFunds(baseChain, baseToken, chain, token);
    }
    
    /**
     * 平台注册token，平台可以使用的通证需要注册。
     * 通证注册之后，平台才能够识别和校验交易。
     * 原则上基础通证(base token)必须经过审核；承兑通证(accept token)可以随意注册。
     * @param request
     * @param chain
     * @param address
     * @return
     */
    @ApiOperation(value="Register token",
    		notes="注册平台允许使用的通证，平台可以使用的通证需要注册。(通证注册之后，平台才能够识别和共识交易)",
    		tags="Register")
    @RequestMapping(value = "/tokens/{chain}/{tokenContractAddress}",
    		method = RequestMethod.POST,
    		produces="application/json;charset=UTF-8")
    @ResponseBody
    public String registerToken(HttpServletRequest request,
    		@PathVariable("chain") String chain,
    		@PathVariable("tokenContractAddress") String address) {
        return fabric.RegisterToken(chain, address);
    }
    
    @ApiOperation(value="Query")
    @RequestMapping(value = "/{id}", method = RequestMethod.GET)
    public JsonResult<Quotation> get(HttpServletRequest request, @PathVariable("id") Integer id) {
    	return new JsonResult<>(HttpStatus.NOT_FOUND);
    }

    @ApiOperation(value="Paging funds")
    @RequestMapping(value = "",
    		method = RequestMethod.GET,
    		produces="application/json;charset=UTF-8")
    @ResponseBody
    public String pageFunds(HttpServletRequest request,
                                            @RequestParam(value = "pageNum") Integer pageNum,
                                            @RequestParam(value = "pageSize") Integer pageSize) {
    	return fabric.ChaincodeQuery("orgchannel", "wallet",
				"{\"Func\":\"query\", \"Args\":[\"funds\", \"page\""
				+ ", \"0x"+Integer.toHexString(pageNum)
				+"\", \"0x"+Integer.toHexString(pageSize)+"\"]}");
    }
    
    @ApiOperation(value="Paging token")
    @RequestMapping(value = "/tokens",
    		method = RequestMethod.GET,
    		produces="application/json;charset=UTF-8")
    @ResponseBody
    public String pageToken(HttpServletRequest request,
                                            @RequestParam(value = "pageNum") Integer pageNum,
                                            @RequestParam(value = "pageSize") Integer pageSize) {
    	return fabric.ChaincodeQuery("orgchannel", "wallet",
				"{\"Func\":\"query\", \"Args\":[\"token\", \"page\""
				+ ", \"0x"+Integer.toHexString(pageNum)
				+"\", \"0x"+Integer.toHexString(pageSize)+"\"]}");
    }
}
