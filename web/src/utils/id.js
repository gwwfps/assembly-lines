import uuid from 'uuid/v4';

const ID_KEY = 'assembly-lines-player-id';

export function getId() {
  const savedId = localStorage.getItem(ID_KEY);
  if (savedId) {
    return savedId;
  }
  const id = uuid();
  localStorage.setItem(ID_KEY, id);
  return id;
}
