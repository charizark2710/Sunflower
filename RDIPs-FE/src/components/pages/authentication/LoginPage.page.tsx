import { useEffect } from 'react';
import { connect } from 'react-redux';

const LoginPage = ({ dispatch }: any) => {
  useEffect(()=> {
    window.location.replace("/api/login")
  },[])

  return (
    <h1> Redirect to login </h1>
  );
};

export default connect()(LoginPage);
