import { combineReducers } from 'redux';
import postReducer from './postReducer';
import navbarTitle from './page'
import authenticate from './authentication';

const reducers = combineReducers({
	posts: postReducer,
	navbarTitle,
	authenticate
});

const rootReducer = (state: any, action : any) => reducers(state, action);
export default rootReducer;