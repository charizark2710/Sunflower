import { combineReducers } from 'redux';
import postReducer from './postReducer';
import navbarTitle from './page'

const reducers = combineReducers({
	posts: postReducer,
	navbarTitle
});

const rootReducer = (state: any, action : any) => reducers(state, action);
export default rootReducer;