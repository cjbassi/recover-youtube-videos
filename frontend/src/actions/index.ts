import { ActionType } from 'typesafe-actions'

import * as actions from './actions'

export type Action = ActionType<typeof actions>
export * from './actions'
export * from './thunks'
