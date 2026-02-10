import { useState, useEffect } from 'react';
import { MapPin, Plus, Pencil, Trash2, Star, Loader2 } from 'lucide-react';
import { useApi } from '../lib/api';
import type { Address, AddressInput } from '../lib/api';
import AddressFormModal from '../components/AddressFormModal';

export default function AddressesPage() {
  const api = useApi();

  const [addresses, setAddresses] = useState<Address[]>([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [editingAddress, setEditingAddress] = useState<Address | null>(null);
  const [deletingId, setDeletingId] = useState<string | null>(null);
  const [settingDefaultId, setSettingDefaultId] = useState<string | null>(null);

  useEffect(() => {
    loadAddresses();
  }, []);

  const loadAddresses = async () => {
    setLoading(true);
    try {
      const response = await api.getAddresses();
      setAddresses(response.data || []);
    } catch (err) {
      console.error('Failed to load addresses:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = async (data: AddressInput) => {
    await api.createAddress(data);
    await loadAddresses();
  };

  const handleUpdate = async (data: AddressInput) => {
    if (!editingAddress) return;
    await api.updateAddress(editingAddress.id, data);
    await loadAddresses();
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this address?')) return;
    setDeletingId(id);
    try {
      await api.deleteAddress(id);
      setAddresses((prev) => prev.filter((a) => a.id !== id));
    } catch (err: any) {
      alert(err.message || 'Failed to delete address');
    } finally {
      setDeletingId(null);
    }
  };

  const handleSetDefault = async (id: string) => {
    setSettingDefaultId(id);
    try {
      await api.setDefaultAddress(id);
      await loadAddresses();
    } catch (err: any) {
      alert(err.message || 'Failed to set default address');
    } finally {
      setSettingDefaultId(null);
    }
  };

  const openEditModal = (addr: Address) => {
    setEditingAddress(addr);
    setShowModal(true);
  };

  const openCreateModal = () => {
    setEditingAddress(null);
    setShowModal(true);
  };

  const closeModal = () => {
    setShowModal(false);
    setEditingAddress(null);
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <Loader2 className="w-8 h-8 animate-spin text-green-600" />
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-3xl mx-auto px-4">
        {/* Header */}
        <div className="flex items-center justify-between mb-8">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">My Addresses</h1>
            <p className="text-gray-600 mt-1">Manage your shipping addresses</p>
          </div>
          <button
            onClick={openCreateModal}
            className="flex items-center gap-2 bg-green-600 text-white px-4 py-2.5 rounded-lg hover:bg-green-700 transition-colors font-medium"
          >
            <Plus className="w-4 h-4" />
            Add Address
          </button>
        </div>

        {/* Address List */}
        {addresses.length === 0 ? (
          <div className="bg-white rounded-lg border border-gray-200 p-12 text-center">
            <MapPin className="w-12 h-12 text-gray-300 mx-auto mb-4" />
            <h3 className="text-lg font-semibold text-gray-900 mb-2">No addresses yet</h3>
            <p className="text-gray-500 mb-6">
              Add a shipping address to speed up your checkout experience.
            </p>
            <button
              onClick={openCreateModal}
              className="inline-flex items-center gap-2 bg-green-600 text-white px-5 py-2.5 rounded-lg hover:bg-green-700 transition-colors font-medium"
            >
              <Plus className="w-4 h-4" />
              Add Your First Address
            </button>
          </div>
        ) : (
          <div className="space-y-4">
            {addresses.map((addr) => (
              <div
                key={addr.id}
                className={`bg-white rounded-lg border p-5 transition-colors ${
                  addr.is_default ? 'border-green-300 ring-1 ring-green-200' : 'border-gray-200'
                }`}
              >
                <div className="flex items-start justify-between gap-4">
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 mb-1">
                      <h3 className="font-semibold text-gray-900">{addr.full_name}</h3>
                      {addr.is_default && (
                        <span className="inline-flex items-center gap-1 text-xs bg-green-100 text-green-700 px-2 py-0.5 rounded-full font-medium">
                          <Star className="w-3 h-3" />
                          Default
                        </span>
                      )}
                    </div>
                    <p className="text-sm text-gray-600">{addr.phone}</p>
                    <p className="text-sm text-gray-600 mt-1">
                      {addr.address_line1}
                      {addr.address_line2 ? `, ${addr.address_line2}` : ''}
                    </p>
                    <p className="text-sm text-gray-600">
                      {addr.city}
                      {addr.state ? `, ${addr.state}` : ''}
                      {addr.postal_code ? ` ${addr.postal_code}` : ''}
                      {' — '}
                      {addr.country}
                    </p>
                  </div>

                  {/* Actions */}
                  <div className="flex items-center gap-2 flex-shrink-0">
                    {!addr.is_default && (
                      <button
                        onClick={() => handleSetDefault(addr.id)}
                        disabled={settingDefaultId === addr.id}
                        className="text-sm text-green-600 hover:text-green-700 font-medium px-3 py-1.5 rounded-md hover:bg-green-50 transition-colors disabled:opacity-50"
                      >
                        {settingDefaultId === addr.id ? 'Setting...' : 'Set Default'}
                      </button>
                    )}
                    <button
                      onClick={() => openEditModal(addr)}
                      className="p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-md transition-colors"
                      title="Edit"
                    >
                      <Pencil className="w-4 h-4" />
                    </button>
                    <button
                      onClick={() => handleDelete(addr.id)}
                      disabled={deletingId === addr.id}
                      className="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-md transition-colors disabled:opacity-50"
                      title="Delete"
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Address Form Modal */}
      <AddressFormModal
        isOpen={showModal}
        onClose={closeModal}
        onSave={editingAddress ? handleUpdate : handleCreate}
        address={editingAddress}
      />
    </div>
  );
}
