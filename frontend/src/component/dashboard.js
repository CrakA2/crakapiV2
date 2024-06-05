
import { Tabs } from "../ui/tabs";
import { Endpoint } from "./endpoint";

export function Dashboard(props) {
  const data = props.data;
  const tabs = [
    {
      title: "RAW Data",
      value: "raw",
      content: (
        <div className="w-full overflow-hidden relative h-full rounded-2xl p-10 text-xl md:text-4xl font-bold text-white bg-gradient-to-br from-purple-700 to-violet-900">
          <Endpoint data={data} api="" />
        </div>
      ),
    },
    {
      title: "StreamLabs",
      value: "StreamLabs",
      content: (
        <div className="w-full overflow-hidden relative h-full rounded-2xl p-10 text-xl md:text-4xl font-bold text-white bg-gradient-to-br from-purple-700 to-violet-900">
          <Endpoint data={data} api="streamlabs"/>

        </div>
      ),
    },
    {
      title: "StreamElements",
      value: "StreamElements",
      content: (
        <div className="w-full overflow-hidden relative h-full rounded-2xl p-10 text-xl md:text-4xl font-bold text-white bg-gradient-to-br from-purple-700 to-violet-900">
          <Endpoint data={data} api="streamelements"/>

        </div>
      ),
    },
    {
      title: "Nightbot",
      value: "NIghtbot",
      content: (
        <div className="w-full overflow-hidden relative h-full rounded-2xl p-10 text-xl md:text-4xl font-bold text-white bg-gradient-to-br from-purple-700 to-violet-900">
          <Endpoint data={data} api="nightbot" />

        </div>
      ),
    },

  ];

  return (
    <div className="h-[10rem] md:h-[40rem] [perspective:1000px] relative b flex flex-col max-w-5xl mx-auto w-full  items-start justify-start my-40">
      <Tabs tabs={tabs} />
    </div>
  );
}
