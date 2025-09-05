"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from "@/components/ui/dialog"

interface AddEntryDialogProps {
  isOpen: boolean
  onClose: () => void
}

export function AddEntryDialog({ isOpen, onClose }: AddEntryDialogProps) {
  const [newEntry, setNewEntry] = useState({ weight: "", date: "", notes: "" })

  const handleAddEntry = () => {
    // Mock add entry logic
    alert(`Added new entry: ${newEntry.weight} lbs on ${newEntry.date}`)
    setNewEntry({ weight: "", date: "", notes: "" })
    onClose()
  }

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add Progress Entry</DialogTitle>
          <DialogDescription>Record your latest measurements</DialogDescription>
        </DialogHeader>
        <div className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="weight">Weight (lbs)</Label>
            <Input
              id="weight"
              type="number"
              placeholder="Enter weight"
              value={newEntry.weight}
              onChange={(e) => setNewEntry({ ...newEntry, weight: e.target.value })}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="date">Date</Label>
            <Input
              id="date"
              type="date"
              value={newEntry.date}
              onChange={(e) => setNewEntry({ ...newEntry, date: e.target.value })}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="notes">Notes (Optional)</Label>
            <Input
              id="notes"
              placeholder="Any notes about your progress..."
              value={newEntry.notes}
              onChange={(e) => setNewEntry({ ...newEntry, notes: e.target.value })}
            />
          </div>
          <Button onClick={handleAddEntry} className="w-full">
            Add Entry
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  )
}
