export default ({ id, players, enabled }) => (
  <div className="column col-6">
    <div className="card">
      <div className="card-header">
        <div className="card-title h5">{id}</div>
        <div className="card-subtitle text-gray">Waiting to start...</div>
      </div>
      <div className="card-body">
        Players ({players.length}
        ):{' '}
        {players.map(player => (
          <span className="chip">{player}</span>
        ))}
      </div>
      <div className="card-footer">
        <button className="btn btn-primary join" value={id} disabled={!enabled}>
          Join
        </button>
      </div>
    </div>
  </div>
);
