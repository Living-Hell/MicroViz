import { useEffect, useState } from "react";
import ReactFlow, {
  MiniMap,
  Controls,
  Background,
  Edge,
  Node,
  BackgroundVariant,
} from "reactflow";
import "reactflow/dist/style.css";
import axios from "axios";
import { Input } from "./components/ui/input";
import { Card, CardContent } from "./components/ui/card";

export default function MicroViz() {
  const [nodes, setNodes] = useState<Node[]>([]);
  const [edges, setEdges] = useState<Edge[]>([]);
  const [search, setSearch] = useState("");

  useEffect(() => {
    axios.get("http://localhost:8080/api/dependencies").then((response) => {
      console.log("data->", response.data);
      const data = response.data;
      const newNodes: Node[] = [];
      const newEdges: Edge[] = [];

      data.forEach((dep: any, index: number) => {
        const sourceNode = {
          id: dep.service_1,
          data: { label: dep.service_1 },
          position: { x: index * 150, y: 100 },
        };
        const targetNode = {
          id: dep.service_2,
          data: { label: dep.service_2 },
          position: { x: index * 150, y: 300 },
        };

        if (!newNodes.find((n) => n.id === dep.service_1))
          newNodes.push(sourceNode);
        if (!newNodes.find((n) => n.id === dep.service_2))
          newNodes.push(targetNode);

        newEdges.push({
          id: `e${dep.service_1}-${dep.service_2}`,
          source: dep.service_1,
          target: dep.service_2,
          label: dep.method,
        });
      });

      setNodes(newNodes);
      setEdges(newEdges);
    });
  }, []);

  const filteredNodes = nodes.filter((node) =>
    node.data.label.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div className="h-screen w-full flex flex-col">
      <header className="p-4 bg-gray-800 text-white text-xl font-bold">
        MicroViz - Dependency Graph
      </header>
      <div className="flex p-4">
        <Input
          className="w-full"
          placeholder="Search services..."
          value={search}
          onChange={(e: any) => setSearch(e.target.value)}
        />
      </div>
      <div className="flex flex-grow">
        <div className="w-1/4 p-4 bg-gray-100">
          <h3 className="font-bold mb-2">Services</h3>
          {filteredNodes.map((node) => (
            <Card
              key={node.id}
              className="mb-2 p-2 cursor-pointer hover:bg-gray-200"
            >
              <CardContent>{node.data.label}</CardContent>
            </Card>
          ))}
        </div>
        <div className="flex-grow h-full">
          <ReactFlow nodes={nodes} edges={edges} fitView>
            <MiniMap />
            <Controls />
            <Background variant={BackgroundVariant.Dots} />
          </ReactFlow>
        </div>
      </div>
    </div>
  );
}
