import { combineReducers } from 'redux';
import postReducer from './postReducer';

const reducers = combineReducers({
	posts: postReducer,
});

const rootReducer = (state: any, action : any) => reducers(state, action);
export default rootReducer;