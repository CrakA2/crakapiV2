import React, { useState } from "react";
import { Input } from "../ui/input";
import { cn } from "./utils/cn.js";
import { Dashboard } from "./dashboard.js";
import { FlipWords } from "../ui/flip-words";

export function SignupFormDemo() {
  const [riotId, setRiotId] = useState("");
  const [tagline, setTagline] = useState("");
  const [data, setData] = useState(null);

  const words = [
    "Ranked Rating",
    "Leaderboard Position",
    "Headshot & KDA",
    "Match History",
  ];
  const bots = ["Nightbot", "OBS", "Streamlabs", "StreamElements"];

  const handleSubmit = async (e) => {
    e.preventDefault();
    const url = `/v1/all/${riotId}/${tagline}`;

    try {
      const response = await fetch(url);
      if (!response.ok) {
        throw new Error("Network response was not ok");
      }
      const data = await response.json();
      setData(data);
      console.log("Form submitted", data);
    } catch (error) {
      alert("Name or Tagline is incorrect");
      console.error(
        "There has been a problem with your fetch operation:",
        error
      );
    }
  };
  if (data) {
    return <Dashboard data={data} />;
  }

  return (
    <div className="flex flex-col items-center justify-center gap-20">
      <div className="flex h-full items-center justify-center gap-20">
        <div className="h-[40rem] flex justify-center items-center px-4">
          <div className="text-4xl mx-auto font-normal text-neutral-600 dark:text-neutral-400">
            GET
            <FlipWords words={words} /> <br />
            With LIVE API ENDPOINTS
          </div>
        </div>
        <div className="max-w-md w-100 h-max mx-auto rounded-none md:rounded-2xl shadow-input pt-10 pb-10 p-12 bg-white dark:bg-black">
          <form className="my-0" onSubmit={handleSubmit}>
            <LabelInputContainer className="mb-4">
              <Input
                id="riotId"
                placeholder="Your Valorant ID here"
                type="text"
                value={riotId}
                onChange={(e) => setRiotId(e.target.value)}
              />
            </LabelInputContainer>
            <LabelInputContainer className="mb-4">
              <Input
                id="tagline"
                placeholder="Your Tag here"
                type="text"
                value={tagline}
                onChange={(e) => setTagline(e.target.value)}
              />
            </LabelInputContainer>

            <button
              className="bg-gradient-to-br relative group/btn from-black dark:from-zinc-900 dark:to-zinc-900 to-neutral-600 block dark:bg-zinc-800 w-full text-white rounded-md h-10 font-medium shadow-[0px_1px_0px_0px_#ffffff40_inset,0px_-1px_0px_0px_#ffffff40_inset] dark:shadow-[0px_1px_0px_0px_var(--zinc-800)_inset,0px_-1px_0px_0px_var(--zinc-800)_inset]"
              type="submit"
            >
              Get Your Stats &rarr;
              <BottomGradient />
            </button>
          </form>
        </div>
      </div>
      <div className="flex justify-center items-center gap-20">
        <div className="text-4xl mx-auto font-normal text-neutral-600 dark:text-neutral-400 mb-20">
          Use with
          <FlipWords words={bots} />
        </div>
      </div>
    </div>
  );
}

const BottomGradient = () => {
  return (
    <>
      <span className="group-hover/btn:opacity-100 block transition duration-500 opacity-0 absolute h-px w-full -bottom-px inset-x-0 bg-gradient-to-r from-transparent via-cyan-500 to-transparent" />
      <span className="group-hover/btn:opacity-100 blur-sm block transition duration-500 opacity-0 absolute h-px w-1/2 mx-auto -bottom-px inset-x-10 bg-gradient-to-r from-transparent via-indigo-500 to-transparent" />
    </>
  );
};

const LabelInputContainer = ({ children, className }) => {
  return (
    <div className={cn("flex flex-col space-y-2 w-full", className)}>
      {children}
    </div>
  );
};
