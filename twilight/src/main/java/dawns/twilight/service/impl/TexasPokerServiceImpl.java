package dawns.twilight.service.impl;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import dawns.twilight.common.annotation.BaseService;
import dawns.twilight.common.base.BaseServiceImpl;
import dawns.twilight.dao.mapper.TexasPokerMapper;
import dawns.twilight.dao.model.TexasPoker;
import dawns.twilight.dao.model.TexasPokerExample;
import dawns.twilight.service.TexasPokerService;

/**
* TexasPokerService实现
* Auto generate on 2019/2/28.
*/
@Service
@Transactional
@BaseService
public class TexasPokerServiceImpl extends BaseServiceImpl<TexasPokerMapper, TexasPoker, TexasPokerExample> implements TexasPokerService {
    @Autowired
    TexasPokerMapper texasPokerMapper;

    @Override
    protected TexasPokerMapper getMapper(){
        return texasPokerMapper;
    }
}
