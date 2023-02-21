import * as React from "react";
import MenuItem from "@mui/material/MenuItem";
import FormControl from "@mui/material/FormControl";
import Select from "@mui/material/Select";
import InputLabel from "@mui/material/InputLabel";
import { useQuery } from "react-query";
import { url } from "../../utils";

const HotfixSelector = ({ hotfixProb, onChange }) => {
  const [hotfix, setHotfix] = React.useState(hotfixProb || "none");
  const { isLoading, error, data } = useQuery("hotfixes", () => {
    return fetch(url("hotfix")).then(async (res) => {
      return await res.json();
    });
  });

  if (isLoading) {
    return <p>Loading...</p>;
  }

  if (error) {
    return <p>Error: {error.message}</p>;
  }

  const hotfixs = (data.data || []).map((hotfix) => hotfix.name);

  const handleChange = (event) => {
    setHotfix(event.target.value);
    onChange(event.target.value);
  };

  return (
    <>
      <FormControl variant="standard" sx={{ m: 0, minWidth: 240 }}>
        <InputLabel>Hotfix</InputLabel>
        <Select
          value={hotfix}
          onChange={handleChange}
          // displayEsmpty
          // inputProps={{ "aria-label": "Without label" }}
          label="Hotfix"
        >
          <MenuItem value="none">
            <em>none</em>
          </MenuItem>
          {hotfixs.map((item) => (
            <MenuItem value={item}>{item}</MenuItem>
          ))}
        </Select>
      </FormControl>
    </>
  );
};

export default HotfixSelector;
