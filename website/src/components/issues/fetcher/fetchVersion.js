import { url, map2Query } from "../../../utils";

export function fetchVersion() {
  return fetch(url("version/maintained")).then(async (res) => {
    const data = await res.json();
    let { data: versions } = data;
    versions.sort();
    versions.reverse();
    return versions || [];
  });
}

export async function fetchVersionByOption({ page = 1, perPage = 100, option = {} }) {
  var queryString = map2Query(option)
  if (queryString !== undefined && queryString !== "") {
    queryString = "&" + queryString
  }

  console.log(
    "fetchVersionWithOption",
    url(`version?page=${page}&per_page=${perPage}${queryString}`)
  );
  try {
    const res = await fetch(url(`version?page=${page}&per_page=${perPage}&${queryString}`));
    const data = await res.json();
    return await data;
  } catch (e) {
    console.log(e);
  }
}
