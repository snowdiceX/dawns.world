package dawns.twilight.controller;

import dawns.twilight.common.base.BaseRestController;
import dawns.twilight.common.web.JsonResult;
import dawns.twilight.dao.model.PublicRandom;
import dawns.twilight.dao.model.PublicRandomExample;
import dawns.twilight.service.PublicRandomService;
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
@RequestMapping("/random")
@Api(value = "PublicRandom api", description = "PublicRandom api",tags = {"PublicRandom"})
public class PublicRandomRestController extends BaseRestController{
    @Autowired
    private PublicRandomService publicRandomService;

    @ApiOperation(value="新增PublicRandom",tags = {"PublicRandom"})
    @RequestMapping(value = "", method = RequestMethod.POST)
    public JsonResult<Integer> add(HttpServletRequest request, @RequestBody PublicRandom publicRandom) {
        return new JsonResult<>(publicRandomService.insert(publicRandom));
    }

    @ApiOperation(value="根据id更新PublicRandom",tags = {"PublicRandom"})
    @RequestMapping(value = "", method = RequestMethod.PUT)
    public JsonResult<Integer> update(HttpServletRequest request, @RequestBody PublicRandom publicRandom) {
        return new JsonResult<>(publicRandomService.updateByPrimaryKey(publicRandom));
    }

    @ApiOperation(value="根据id查询PublicRandom",tags = {"PublicRandom"})
    @RequestMapping(value = "/{id}", method = RequestMethod.GET)
    public JsonResult<PublicRandom> get(HttpServletRequest request, @PathVariable("id") Integer id) {
        PublicRandom publicRandom=publicRandomService.selectByPrimaryKey(id);
        if(publicRandom!=null){
            return new JsonResult<>(publicRandom);
        }else{
            return new JsonResult<>(HttpStatus.NOT_FOUND);
        }
    }

    @ApiOperation(value="根据id删除PublicRandom",tags = {"PublicRandom"})
    @RequestMapping(value = "/{id}", method = RequestMethod.DELETE)
    public JsonResult<Integer> delete(HttpServletRequest request, @PathVariable("id") Integer id) {
        return new JsonResult<>(publicRandomService.deleteByPrimaryKey(id));
    }

    @ApiOperation(value="分页查询PublicRandom",tags = {"PublicRandom"})
    @RequestMapping(value = "", method = RequestMethod.GET)
    public JsonResult<List<PublicRandom>> page(HttpServletRequest request,
                                            @RequestParam(value = "pageNum") Integer pageNum,
                                            @RequestParam(value = "pageSize") Integer pageSize) {
        PublicRandomExample publicRandomExample=new PublicRandomExample();
        publicRandomExample.setOrderByClause("id");
        return new JsonResult<>(publicRandomService.selectByExampleForStartPage(publicRandomExample, pageNum,pageSize));
    }
}
