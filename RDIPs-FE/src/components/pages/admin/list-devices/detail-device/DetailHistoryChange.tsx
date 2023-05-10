import { useLocation } from 'react-router';

const DetailHistoryChange = () => {
  let { state } = useLocation();
  const detailHistoryChange: any = state;
  console.log(detailHistoryChange);

  return (
    <div style={{ padding: '0 30px', backgroundColor: 'white', minHeight: '80vh' }}>
      <section>
        <h3>Device History Change</h3>
      </section>
    </div>
  );
};

export default DetailHistoryChange;
