// @flow
import * as Avatar from '../desktop/remote/sync-avatar-props.desktop'

export const serialize = {
  ...Avatar.serialize,
  badgeKeys: v => v.keySeq().toArray(),
  badgeMap: (v, o) =>
    v
      .keySeq()
      .toArray()
      .reduce((map, k) => {
        if (!o || v.get(k) !== o.get(k)) {
          map[k] = v.get(k)
        }
        return map
      }, {}),
  broken: v => v,
  conversations: v => v,
  externalRemoteWindow: v => v,
  fileRows: v => v,
  isAsyncWriteHappening: v => v,
  loggedIn: v => v,
  username: v => v,
  windowComponent: v => v,
  windowOpts: v => v,
  windowParam: v => v,
  windowTitle: v => v,
}

const initialState = {
  badgeInfo: {},
  badgeMap: {},
  config: {},
}
export const deserialize = (state: any = initialState, props: any) => {
  if (!props) return state
  // We always add to the map
  const badgeMap = {
    ...state.badgeMap,
    ...(props.badgeMap || {}),
  }

  const badgeInfo = (props.badgeKeys || []).reduce((map, k) => {
    map[k] = badgeMap[k]
    return map
  }, {})

  const newState = {
    ...state,
    ...props,
    badgeInfo,
    badgeMap,
  }
  return Avatar.deserialize(newState, props)
}