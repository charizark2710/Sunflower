import { useLocation } from 'react-router';

const DetailHistoryLog = () => {
  let { state } = useLocation();
  // const detailHistoryLog: any = state;
  // console.log(detailHistoryLog);

  return (
    <div style={{ padding: '0 30px', backgroundColor: 'white', minHeight: '80vh' }}>
      <section>
        <h3>Device History Log</h3>
      </section>
    </div>
  );
};

export default DetailHistoryLog;
