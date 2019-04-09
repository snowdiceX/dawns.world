package dawns.twilight.controller;

import java.util.List;

import javax.servlet.http.HttpServletRequest;

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
import dawns.twilight.common.chain.FabricService;
import dawns.twilight.common.web.JsonResult;
import dawns.twilight.common.web.ResponseToken;
import dawns.twilight.common.web.ResponseWallet;
import dawns.twilight.dao.model.Quotation;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;

@Controller
@RequestMapping("/funds")
@Api(value = "Funds api", description = "Funds api",tags = {"3. Funds"})
public class FundsRestController extends BaseRestController{
    
    @Autowired
    private FabricService fabric;

    @ApiOperation(value="Deposit")
    @RequestMapping(value = "", method = RequestMethod.POST)
    public JsonResult<Integer> add(HttpServletRequest request, @RequestBody Quotation quotation) {
        return new JsonResult<>(null);
    }

    @ApiOperation(value="Withdraw")
    @RequestMapping(value = "/{id}", method = RequestMethod.DELETE)
    public JsonResult<Integer> delete(HttpServletRequest request, @PathVariable("id") Integer id) {
        return new JsonResult<>(null);
    }
    
    @ApiOperation(value="Import token")
    @RequestMapping(value = "/tokens/{chain}/{tokenContractAddress}", method = RequestMethod.POST)
    public JsonResult<ResponseToken> importToken(HttpServletRequest request,
    		@PathVariable("chain") String chain,
    		@PathVariable("tokenContractAddress") String address) {
    	JsonResult<ResponseToken> result = new JsonResult<>(HttpStatus.OK);
    	String ret = fabric.ImportToken(chain, address);
    	JSONObject obj = JSONObject.parseObject(ret);
    	result.setCode(obj.getInteger("code"));
    	result.setMessage(obj.getString("message"));
    	ResponseToken token = new ResponseToken();
    	obj = obj.getJSONObject("result");
    	token.setNetwork(obj.getString("network"));
    	token.setAddress(obj.getString("address"));
//    	token.setHeight(obj.getString("height"));
//    	token.setTxid(obj.getString("txid"));
    	token.setName(obj.getString("name"));
    	token.setSymbol(obj.getString("symbol"));
    	token.setDecimals(obj.getInteger("decimals"));
    	result.setData(token);
        return result;
    }
    
    @ApiOperation(value="Query")
    @RequestMapping(value = "/{id}", method = RequestMethod.GET)
    public JsonResult<Quotation> get(HttpServletRequest request, @PathVariable("id") Integer id) {
    	return new JsonResult<>(HttpStatus.NOT_FOUND);
    }

    @ApiOperation(value="Paging")
    @RequestMapping(value = "", method = RequestMethod.GET)
    public JsonResult<List<Quotation>> page(HttpServletRequest request,
                                            @RequestParam(value = "pageNum") Integer pageNum,
                                            @RequestParam(value = "pageSize") Integer pageSize) {
        return new JsonResult<>(null);
    }
}
