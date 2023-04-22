import { useNavigate } from 'react-router-dom';
import './ListUsers.scss';

const ListUsers = () => {
  const navigate = useNavigate();
  // function navigateToDetailPage(detail: any) {
  //   navigate('/detail-user', { replace: false, state: detail });
  // }

  return (
    <div className='list-users-container'>
     Users List here
    </div>
  );
};

export default ListUsers;