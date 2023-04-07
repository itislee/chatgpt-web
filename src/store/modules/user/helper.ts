import { ss } from '@/utils/storage'

const LOCAL_NAME = 'userStorage'

export interface UserInfo {
  avatar: string
  name: string
  description: string
}

export interface UserState {
  userInfo: UserInfo
}

export function getCookie(param: string): string {
  if (document.cookie.length > 0) {
    const list = document.cookie.split('; ')
    for (let i = 0; i < list.length; i++) {
      const arr = list[i].split('=')
      if (arr[0] === param)
        return arr[1]
    }
    return ''
  }
}

export function defaultSetting(): UserState {
  return {
    userInfo: {
      avatar: 'https://q.qlogo.cn/qzapp/101570536/'.concat(getCookie('open_id'), '/100'),
      name: getCookie('nick'),
      description: 'Star on <a href="https://github.com/Chanzhaoyu/chatgpt-bot" class="text-blue-500" target="_blank" >Github</a>',
    },
  }
}

export function getLocalState(): UserState {
  const localSetting: UserState | undefined = ss.get(LOCAL_NAME)
  return { ...defaultSetting(), ...localSetting }
}

export function setLocalState(setting: UserState): void {
  ss.set(LOCAL_NAME, setting)
}
