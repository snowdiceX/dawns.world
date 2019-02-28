package dawns.twilight.service.impl;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import dawns.twilight.common.annotation.BaseService;
import dawns.twilight.common.base.BaseServiceImpl;
import dawns.twilight.dao.mapper.GameMapper;
import dawns.twilight.dao.model.Game;
import dawns.twilight.dao.model.GameExample;
import dawns.twilight.service.GameService;

/**
* GameService实现
* Auto generate on 2019/2/28.
*/
@Service
@Transactional
@BaseService
public class GameServiceImpl extends BaseServiceImpl<GameMapper, Game, GameExample> implements GameService {
    @Autowired
    GameMapper gameMapper;

    @Override
    protected GameMapper getMapper(){
        return gameMapper;
    }
}
