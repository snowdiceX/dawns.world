package dawns.twilight.controller;

import dawns.twilight.common.base.BaseRestController;
import dawns.twilight.common.web.JsonResult;
import dawns.twilight.dao.model.Quotation;
import dawns.twilight.dao.model.QuotationExample;
import dawns.twilight.service.QuotationService;
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
@RequestMapping("/funds")
@Api(value = "Funds api", description = "Funds api",tags = {"3. Funds"})
public class FundsRestController extends BaseRestController{
    @Autowired
    private QuotationService quotationService;

    @ApiOperation(value="Deposit")
    @RequestMapping(value = "", method = RequestMethod.POST)
    public JsonResult<Integer> add(HttpServletRequest request, @RequestBody Quotation quotation) {
        return new JsonResult<>(quotationService.insert(quotation));
    }

    @ApiOperation(value="Withdraw")
    @RequestMapping(value = "/{id}", method = RequestMethod.DELETE)
    public JsonResult<Integer> delete(HttpServletRequest request, @PathVariable("id") Integer id) {
        return new JsonResult<>(quotationService.deleteByPrimaryKey(id));
    }
    
    @ApiOperation(value="Query")
    @RequestMapping(value = "/{id}", method = RequestMethod.GET)
    public JsonResult<Quotation> get(HttpServletRequest request, @PathVariable("id") Integer id) {
        Quotation quotation=quotationService.selectByPrimaryKey(id);
        if(quotation!=null){
            return new JsonResult<>(quotation);
        }else{
            return new JsonResult<>(HttpStatus.NOT_FOUND);
        }
    }

    @ApiOperation(value="Paging")
    @RequestMapping(value = "", method = RequestMethod.GET)
    public JsonResult<List<Quotation>> page(HttpServletRequest request,
                                            @RequestParam(value = "pageNum") Integer pageNum,
                                            @RequestParam(value = "pageSize") Integer pageSize) {
        QuotationExample quotationExample=new QuotationExample();
        quotationExample.setOrderByClause("id");
        return new JsonResult<>(quotationService.selectByExampleForStartPage(quotationExample, pageNum,pageSize));
    }
}
