import { url, map2Query } from "../../../utils";

export async function fetchHotfix({ page = 1, perPage = 100, versionName }) {
  // Judge the condition here cause the useQuery function must on top level of Parent function.
  if (versionName == "none") {
    return {}
  }

  var versionOption = composeVersionOption(versionName)
  return fetchVersionByOption({ page: page, perPage: perPage, option: versionOption })
    .then(async (data) => {
      let { data: versions } = data;
      versions.sort(
        (a, b) => {
          return a > b ? -1 : 1;
        }
      );
      return versions || []
    });
}


