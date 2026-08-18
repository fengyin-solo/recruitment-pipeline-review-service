# 招聘系统 (Recruit)

纯 Go 标准库实现的招聘管理后端服务。

## 运行说明

```bash
cd 12-recruit/origin
go run ./cmd/server
```

默认监听 `:8080`，可通过环境变量 `PORT` 或 `ADDR` 修改。

## API 列表

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/jobs | 创建职位 |
| GET | /api/jobs | 职位列表（支持 department/status/keyword 筛选 + 分页） |
| GET | /api/jobs/{id} | 获取职位 |
| PUT | /api/jobs/{id} | 更新职位 |
| DELETE | /api/jobs/{id} | 删除职位 |
| POST | /api/jobs/batch-close | 批量关闭职位 |
| POST | /api/candidates | 创建候选人 |
| GET | /api/candidates | 候选人列表（支持 keyword/status 筛选 + 分页） |
| GET | /api/candidates/{id} | 获取候选人 |
| PUT | /api/candidates/{id} | 更新候选人 |
| DELETE | /api/candidates/{id} | 删除候选人 |
| POST | /api/resumes | 创建简历（校验 candidate_id 外键） |
| GET | /api/resumes | 简历列表（支持 candidate_id 筛选 + 分页） |
| GET | /api/resumes/{id} | 获取简历 |
| PUT | /api/resumes/{id} | 更新简历 |
| DELETE | /api/resumes/{id} | 删除简历 |
| POST | /api/interviews | 创建面试（校验 job_id/candidate_id 外键） |
| GET | /api/interviews | 面试列表（支持 job_id/candidate_id/status 筛选 + 分页） |
| GET | /api/interviews/{id} | 获取面试 |
| PUT | /api/interviews/{id} | 更新面试（含状态机校验） |
| DELETE | /api/interviews/{id} | 删除面试 |
| POST | /api/interviews/batch-update-status | 批量更新面试状态 |
| POST | /api/offers | 创建 Offer（校验 job_id/candidate_id 外键） |
| GET | /api/offers | Offer 列表（支持 job_id/candidate_id/status 筛选 + 分页） |
| GET | /api/offers/{id} | 获取 Offer |
| PUT | /api/offers/{id} | 更新 Offer（含状态机校验，accepted 时可能触发职位 filled） |
| DELETE | /api/offers/{id} | 删除 Offer |
| GET | /api/stats/departments | 按部门统计职位数 |
| GET | /api/stats/funnel | 按职位统计候选人/面试/Offer 漏斗 |

## 技术栈

- Go 1.22（使用增强版 ServeMux 方法路由）
- 纯标准库，零第三方依赖
- 内存存储（线程安全）
- 统一响应格式：`{"code":0,"message":"ok","data":...}`
