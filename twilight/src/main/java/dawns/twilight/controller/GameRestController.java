package dawns.twilight.controller;

import dawns.twilight.common.base.BaseRestController;
import dawns.twilight.common.web.JsonResult;
import dawns.twilight.dao.model.Game;
import dawns.twilight.dao.model.GameExample;
import dawns.twilight.service.GameService;
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
@RequestMapping("/game")
@Api(value = "Game api", description = "Game api",tags = {"Game"})
public class GameRestController extends BaseRestController{
    @Autowired
    private GameService gameService;

    @ApiOperation(value="新增Game",tags = {"Game"})
    @RequestMapping(value = "", method = RequestMethod.POST)
    public JsonResult<Integer> add(HttpServletRequest request, @RequestBody Game game) {
        return new JsonResult<>(gameService.insert(game));
    }

    @ApiOperation(value="根据id更新Game",tags = {"Game"})
    @RequestMapping(value = "", method = RequestMethod.PUT)
    public JsonResult<Integer> update(HttpServletRequest request, @RequestBody Game game) {
        return new JsonResult<>(gameService.updateByPrimaryKey(game));
    }

    @ApiOperation(value="根据id查询Game",tags = {"Game"})
    @RequestMapping(value = "/{id}", method = RequestMethod.GET)
    public JsonResult<Game> get(HttpServletRequest request, @PathVariable("id") Integer id) {
        Game game=gameService.selectByPrimaryKey(id);
        if(game!=null){
            return new JsonResult<>(game);
        }else{
            return new JsonResult<>(HttpStatus.NOT_FOUND);
        }
    }

    @ApiOperation(value="根据id删除Game",tags = {"Game"})
    @RequestMapping(value = "/{id}", method = RequestMethod.DELETE)
    public JsonResult<Integer> delete(HttpServletRequest request, @PathVariable("id") Integer id) {
        return new JsonResult<>(gameService.deleteByPrimaryKey(id));
    }

    @ApiOperation(value="分页查询Game",tags = {"Game"})
    @RequestMapping(value = "", method = RequestMethod.GET)
    public JsonResult<List<Game>> page(HttpServletRequest request,
                                            @RequestParam(value = "pageNum") Integer pageNum,
                                            @RequestParam(value = "pageSize") Integer pageSize) {
        GameExample gameExample=new GameExample();
        gameExample.setOrderByClause("id");
        return new JsonResult<>(gameService.selectByExampleForStartPage(gameExample, pageNum,pageSize));
    }
}
