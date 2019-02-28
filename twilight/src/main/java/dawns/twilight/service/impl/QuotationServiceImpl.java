package dawns.twilight.service.impl;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import dawns.twilight.common.annotation.BaseService;
import dawns.twilight.common.base.BaseServiceImpl;
import dawns.twilight.dao.mapper.QuotationMapper;
import dawns.twilight.dao.model.Quotation;
import dawns.twilight.dao.model.QuotationExample;
import dawns.twilight.service.QuotationService;

/**
* QuotationService实现
* Auto generate on 2019/2/28.
*/
@Service
@Transactional
@BaseService
public class QuotationServiceImpl extends BaseServiceImpl<QuotationMapper, Quotation, QuotationExample> implements QuotationService {
    @Autowired
    QuotationMapper quotationMapper;

    @Override
    protected QuotationMapper getMapper(){
        return quotationMapper;
    }
}
