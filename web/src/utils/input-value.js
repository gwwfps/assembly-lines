export default el => el.events('input').map(ev => ev.target.value.trim());
