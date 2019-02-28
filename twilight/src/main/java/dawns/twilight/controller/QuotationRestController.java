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
@RequestMapping("/quotation")
@Api(value = "Quotation api", description = "Quotation api",tags = {"Quotation"})
public class QuotationRestController extends BaseRestController{
    @Autowired
    private QuotationService quotationService;

    @ApiOperation(value="新增Quotation",tags = {"Quotation"})
    @RequestMapping(value = "", method = RequestMethod.POST)
    public JsonResult<Integer> add(HttpServletRequest request, @RequestBody Quotation quotation) {
        return new JsonResult<>(quotationService.insert(quotation));
    }

    @ApiOperation(value="根据id更新Quotation",tags = {"Quotation"})
    @RequestMapping(value = "", method = RequestMethod.PUT)
    public JsonResult<Integer> update(HttpServletRequest request, @RequestBody Quotation quotation) {
        return new JsonResult<>(quotationService.updateByPrimaryKey(quotation));
    }

    @ApiOperation(value="根据id查询Quotation",tags = {"Quotation"})
    @RequestMapping(value = "/{id}", method = RequestMethod.GET)
    public JsonResult<Quotation> get(HttpServletRequest request, @PathVariable("id") Integer id) {
        Quotation quotation=quotationService.selectByPrimaryKey(id);
        if(quotation!=null){
            return new JsonResult<>(quotation);
        }else{
            return new JsonResult<>(HttpStatus.NOT_FOUND);
        }
    }

    @ApiOperation(value="根据id删除Quotation",tags = {"Quotation"})
    @RequestMapping(value = "/{id}", method = RequestMethod.DELETE)
    public JsonResult<Integer> delete(HttpServletRequest request, @PathVariable("id") Integer id) {
        return new JsonResult<>(quotationService.deleteByPrimaryKey(id));
    }

    @ApiOperation(value="分页查询Quotation",tags = {"Quotation"})
    @RequestMapping(value = "", method = RequestMethod.GET)
    public JsonResult<List<Quotation>> page(HttpServletRequest request,
                                            @RequestParam(value = "pageNum") Integer pageNum,
                                            @RequestParam(value = "pageSize") Integer pageSize) {
        QuotationExample quotationExample=new QuotationExample();
        quotationExample.setOrderByClause("id");
        return new JsonResult<>(quotationService.selectByExampleForStartPage(quotationExample, pageNum,pageSize));
    }
}
