package dawns.twilight.controller;

import dawns.twilight.common.base.BaseRestController;
import dawns.twilight.common.web.JsonResult;
import dawns.twilight.dao.model.WalletAddress;
import dawns.twilight.dao.model.WalletAddressExample;
import dawns.twilight.service.WalletAddressService;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RequestParam;

import javax.servlet.http.HttpServletRequest;
import java.util.List;

@Controller
@RequestMapping("/wallet")
@Api(value = "WalletAddress api", description = "WalletAddress api",tags = {"WalletAddress"})
public class WalletAddressRestController extends BaseRestController{
    @Autowired
    private WalletAddressService walletAddressService;

    @ApiOperation(value="新增WalletAddress",tags = {"WalletAddress"})
    @RequestMapping(value = "", method = RequestMethod.POST)
    public JsonResult<Integer> add(HttpServletRequest request, @RequestBody WalletAddress walletAddress) {
        return new JsonResult<>(walletAddressService.insert(walletAddress));
    }

    @ApiOperation(value="根据id更新WalletAddress",tags = {"WalletAddress"})
    @RequestMapping(value = "", method = RequestMethod.PUT)
    public JsonResult<Integer> update(HttpServletRequest request, @RequestBody WalletAddress walletAddress) {
        return new JsonResult<>(walletAddressService.updateByPrimaryKey(walletAddress));
    }

    @ApiOperation(value="根据id查询WalletAddress",tags = {"WalletAddress"})
    @RequestMapping(value = "/{id}", method = RequestMethod.GET)
    public JsonResult<WalletAddress> get(HttpServletRequest request, @PathVariable("id") Integer id) {
        WalletAddress walletAddress=walletAddressService.selectByPrimaryKey(id);
        if(walletAddress!=null){
            return new JsonResult<>(walletAddress);
        }else{
            return new JsonResult<>(HttpStatus.NOT_FOUND);
        }
    }

    @ApiOperation(value="根据id删除WalletAddress",tags = {"WalletAddress"})
    @RequestMapping(value = "/{id}", method = RequestMethod.DELETE)
    public JsonResult<Integer> delete(HttpServletRequest request, @PathVariable("id") Integer id) {
        return new JsonResult<>(walletAddressService.deleteByPrimaryKey(id));
    }

    @ApiOperation(value="分页查询WalletAddress",tags = {"WalletAddress"})
    @RequestMapping(value = "", method = RequestMethod.GET)
    public JsonResult<List<WalletAddress>> page(HttpServletRequest request,
                                            @RequestParam(value = "pageNum") Integer pageNum,
                                            @RequestParam(value = "pageSize") Integer pageSize) {
        WalletAddressExample walletAddressExample=new WalletAddressExample();
        walletAddressExample.setOrderByClause("id");
        return new JsonResult<>(walletAddressService.selectByExampleForStartPage(walletAddressExample, pageNum,pageSize));
    }
}
