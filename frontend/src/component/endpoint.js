import { Textfield } from "../ui/textfield.js";

export function Endpoint(props) {
  const data = props.data;
  const url = window.location.href + "v1";
  const region = data.region;
  const puuid = data.puuid.replace(/-/g, "");

  if (!data) {
    return <div>Loading...</div>;
  }

  if (props.api === "streamlabs") {
    return (
      <div>
        <h3>
          Hi, {data.name} #{data.tag}
        </h3>
        <div className=" text-lg ">
          <p>
            Puuid: {data.puuid} Region: {data.region} !command add !commandname then endpoints
          </p>
          <p>KDA: {data.kda}</p>
          <Textfield link={`$(urlfetch ${url}/kd/${region}/${puuid})`} />
          <p>Headshot</p>
          <Textfield link={`$(urlfetch ${url}/hs/${region}/${puuid})`} />

          <p>
            RR: {data.rr.CurrentTierPatched}, Ranking: {data.rr.RankingInTier}
          </p>
          <Textfield link={`$(urlfetch ${url}/rr/${region}/${puuid})`} />
          <p>
            Leaderboard: {data.leaderboard}
            <span> only visible if you are on the leaderboard </span>
          </p>
          <Textfield link={`$(urlfetch ${url}/lb/${region}/${puuid})`} />
          <p>Win/Loss Record: {data.winloss.Record.join(", ")}</p>
          <Textfield link={`$(urlfetch ${url}/wl/${region}/${puuid})`} />
          <p>
            Wins: {data.wins} &nbsp; Losses: {data.losses} &nbsp; Draws:{" "}
            {data.draws}
          </p>
        </div>
      </div>
    );
  }

  if (props.api === "streamelements") {
    return (
      <div>
        <h3>
          Hi, {data.name} #{data.tag}
        </h3>
        <div className=" text-lg ">
          <p>
            Puuid: {data.puuid} Region: {data.region} Use !command add !commandname then endpoints
          </p>
          <p>KDA: {data.kda}</p>
          <Textfield link={`\${urlfetch ${url}/kd/${region}/${puuid}}`} />
          <p>Headshot</p>
          <Textfield link={`\${urlfetch ${url}/hs/${region}/${puuid}}`} />

          <p>
            RR: {data.rr.CurrentTierPatched}, Ranking: {data.rr.RankingInTier}
          </p>
          <Textfield link={`\${urlfetch ${url}/rr/${region}/${puuid}}`} />
          <p>
            Leaderboard: {data.leaderboard}
            <span> only visible if you are on the leaderboard </span>
          </p>
          <Textfield link={`\${urlfetch ${url}/lb/${region}/${puuid}}`} />
          <p>Win/Loss Record: {data.winloss.Record.join(", ")}</p>
          <Textfield link={`\${urlfetch ${url}/wl/${region}/${puuid}}`} />
          <p>
            Wins: {data.wins} &nbsp; Losses: {data.losses} &nbsp; Draws:{" "}
            {data.draws}
          </p>
        </div>
      </div>
    );
  }

  if (props.api === "nightbot") {
    return (
      <div>
        <h3>
          Hi, {data.name} #{data.tag}
        </h3>
        <div className=" text-lg ">
          <p>
            Puuid: {data.puuid} Region: {data.region} !command add !commandname then endpoints
          </p>
          <p>KDA: {data.kda}</p>
          <Textfield link={`$(urlfetch ${url}/kd/${region}/${puuid})`} />
          <p>Headshot</p>
          <Textfield link={`$(urlfetch ${url}/hs/${region}/${puuid})`} />

          <p>
            RR: {data.rr.CurrentTierPatched}, Ranking: {data.rr.RankingInTier}
          </p>
          <Textfield link={`$(urlfetch ${url}/rr/${region}/${puuid})`} />
          <p>
            Leaderboard: {data.leaderboard}
            <span> only visible if you are on the leaderboard </span>
          </p>
          <Textfield link={`$(urlfetch ${url}/lb/${region}/${puuid})`} />
          <p>Win/Loss Record: {data.winloss.Record.join(", ")}</p>
          <Textfield link={`$(urlfetch ${url}/wl/${region}/${puuid})`} />
          <p>
            Wins: {data.wins} &nbsp; Losses: {data.losses} &nbsp; Draws:{" "}
            {data.draws}
          </p>
        </div>
      </div>
    );
  }

  return (
    <div>
      <h3>
        Hi, {data.name} #{data.tag}
      </h3>
      <div className=" text-lg ">
        <p>
          Puuid: {data.puuid} Region: {data.region}
        </p>
        <p>KDA: {data.kda}</p>
        <Textfield link={`${url}/kd/${region}/${puuid}`} />
        <p>Headshot</p>
        <Textfield link={`${url}/hs/${region}/${puuid}`} />

        <p>
          RR: {data.rr.CurrentTierPatched}, Ranking: {data.rr.RankingInTier}
        </p>
        <Textfield link={`${url}/rr/${region}/${puuid}`} />
        <p>
          Leaderboard: {data.leaderboard}
          <span> only visible if you are on the leaderboard </span>
        </p>
        <Textfield link={`${url}/lb/${region}/${puuid}`} />
        <p>Win/Loss Record: {data.winloss.Record.join(", ")}</p>
        <Textfield link={`${url}/wl/${region}/${puuid}`} />
        <p>
          Wins: {data.wins} &nbsp; Losses: {data.losses} &nbsp; Draws:{" "}
          {data.draws}
        </p>
      </div>
    </div>
  );
}
