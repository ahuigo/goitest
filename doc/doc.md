---
title: Golang：从API单元测试生成集成测试
date: 2024-11-30
---
# 前序
> https://medium.com/p/c315a89e523e/edit
日常API接口开发时，许多时间会花在测试上。测试的方法主要有以下几种：

1. 使用postman(或chrome/curl)访问 dev api接口测试
    1.　只是手动测试测试观察
    2.　每次要重启dev　api
    3.　测试不方便共享
    4.　需要在切换app与ide之间来回切换
2. 裸写go test单元测试
    - 优点:　
        1. 自动化
        1. 一键生成覆盖率
        1. 一键debug 断点测试
    - 缺点:　代码比较繁琐; 不能生成集成测试
3. 使用goitest 工具测试
    1. 自动化单元测试
    2. 一键生成集成测试

# Integrate rules 架构

    req case def：
        name: string // unique
        request 实体
            reqTpl // 模板定义 , req-tpl
        response 实体
        expectRule[]: (expectCase?)
            name: string // unique
    req case save：

    testExpect:

# todo
- []save req+test_rules
- []load req data
