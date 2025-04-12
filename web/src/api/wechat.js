import request from '@/utils/request'

export function getAuthUrl() {
  return request({
    url: '/api/v1/wechat/auth',
    method: 'get'
  })
} 